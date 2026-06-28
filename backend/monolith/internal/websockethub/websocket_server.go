package websockethub

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"

	wsPB "monolith/internal/websockethub/proto"
	datastructures "monolith/pkg/data_structures"
)

type WebsocketServerOptions struct {
	Roles []string `json:"roles"`

	Logger       zerolog.Logger `json:"-"`
	DebugMode    bool           `json:"-"`
	TraceMode    bool           `json:"-"`
	TraceVerbose bool           `json:"-"`
}

type WebsocketServer struct {
	upgrader          *websocket.Upgrader
	msgScheduler      MessageScheduler
	ConnectionsByRole map[string]*ConnectionsByRole
	// roles             []string

	logger    zerolog.Logger
	debugMode bool
	traceMode bool
}

func NewWebsocketServer(config WebsocketServerOptions) *WebsocketServer {
	ws := &WebsocketServer{}

	ws.upgrader = ws.newUpgrader()
	ws.msgScheduler = *NewMessageScheduler(10*time.Millisecond, 500)

	connectionsByRole := make(map[string]*ConnectionsByRole)
	for _, role := range config.Roles {
		connectionsByRole[role] = &ConnectionsByRole{
			activeConnections: &datastructures.SyncMap[uint32, *ConnectionData]{},
			// channels:          make(map[string]ChannelActions),
			channels: make([]ChannelActions, len(wsPB.Channel_name)),
		}
	}
	ws.ConnectionsByRole = connectionsByRole

	logger := config.Logger
	if logger.GetLevel() == zerolog.NoLevel {
		output := zerolog.ConsoleWriter{Out: os.Stdout}
		logger = zerolog.New(output).Level(zerolog.InfoLevel).With().Timestamp().Logger()
	}
	ws.logger = logger

	return ws
}

func (ws *WebsocketServer) Run() {
	ws.msgScheduler.Start()
}

func (ws *WebsocketServer) Stop() {
	ws.msgScheduler.Stop()
}

func (ws *WebsocketServer) RegisterChannel(roles []string, name wsPB.Channel, channel ChannelActions) error {
	// Register for all roles
	isEmptyRoles := len(roles) == 0
	if isEmptyRoles {
		for role := range ws.ConnectionsByRole {
			cbr, _ := ws.ConnectionsByRole[role]
			cbr.channels[name] = channel
		}
		return nil
	}

	for _, role := range roles {
		cbr, ok := ws.ConnectionsByRole[role]
		if !ok {
			return fmt.Errorf("role %s does not exists", role)
		}
		cbr.channels[name] = channel
	}
	return nil
}

// TODO:
// func (ws *WebsocketServer) RegisterRoles(roles []string) {}

func (ws *WebsocketServer) WebsocketAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// TODO: use valkey to get info by [token], fallback to main db if cache miss
		// token := strings.TrimPrefix(authHeader, "Bearer ")

		// TODO: Use local casbin instead of async network service
		// done := make(chan bool, 1)
		// go func() {
		// 	select {
		// 	case ok := <-done:
		// 		if !ok {
		// 			http.Error(w, "unauthorized", http.StatusUnauthorized)
		// 		}
		// 	case <-time.After(2 * time.Second):
		// 		http.Error(w, "validation timeout", http.StatusGatewayTimeout)
		// 	}
		// }()

		// TODO: context
		// ctx := context.WithValue(parent context.Context, key any, val any)
		// next.ServeHTTP(w, r.WithContext(ctx))
		//
		next.ServeHTTP(w, r)
	})
}

func (ws *WebsocketServer) WebsocketUpgradeHandler(ctx *gin.Context) {
	idParam := ctx.Param("id")
	connID64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.String(http.StatusBadRequest, "invalid 'id' parameter: must be a positive integer")
		return
	}
	connID := uint32(connID64)
	// TODO: handle unregistered users
	// if connID == 0 {
	// }
	connRole := ctx.Param("role")

	connPool, ok := ws.ConnectionsByRole[connRole]
	if !ok {
		ctx.JSON(http.StatusBadRequest, "no such role")
		return
	}

	conn, err := ws.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to upgrade connection")
		return
	}

	connData := &ConnectionData{
		ID:             connID,
		ConnectionPool: connPool,
		// subscriptions:  make(map[string][]uint32),
		subscriptions: make([][]uint32, 1),
	}
	conn.SetSession(connData)

	// TODO: change active connections: currently we are using IDs from DB,
	// however when unregistered user tries to log in we need to generate id for him =>
	// => generate id as connection index in pool for every connection (also improves caching)
	connPool.activeConnections.Store(connID, connData)
}

func (ws *WebsocketServer) newUpgrader() *websocket.Upgrader {
	upgrader := websocket.NewUpgrader()

	upgrader.OnOpen(ws.handleWebsocketConnect)
	upgrader.OnClose(ws.handleWebsocketDisconnect)

	// If we care about error of writing back to connection, than retry. Else handling write error is NO-OP
	upgrader.OnMessage(func(conn *websocket.Conn, messageType websocket.MessageType, data []byte) {
		// TODO: change to event
		// ws.logger.Debug().
		// 	Int("bytes received", len(data)).
		// 	Any("from client", conn.RemoteAddr()).
		// 	Send()

		msg := &wsPB.Request{}
		err := proto.Unmarshal(data, msg)
		if err != nil {
			sendErrorResponse(conn, "", http.StatusBadRequest, "incorrect message")
			return
		}
		// if msg.RequestId == "" {
		// }
		resp := wsPB.ResponseMessage{
			RequestId: msg.RequestId,
		}

		switch payload := msg.Payload.(type) {
		case *wsPB.Request_Publish:
			err = ws.handlePublish(msg.Channel, payload.Publish.Data, conn)
			if err != nil {
				resp.StatusCode = http.StatusInternalServerError
				resp.ErrorMessage = err.Error()
			} else {
				resp.StatusCode = http.StatusOK
			}
		case *wsPB.Request_Subscribe:
			data, err := ws.handleSubscribe(msg.Channel, payload.Subscribe.Indexes, conn)
			if err != nil {
				resp.StatusCode = http.StatusInternalServerError
				resp.ErrorMessage = err.Error()
			} else {
				resp.StatusCode = http.StatusOK
				resp.Data = data
			}
		case *wsPB.Request_Unsubscribe:
			ws.handleUnsubscribe(msg.Channel, payload.Unsubscribe.Indexes, conn)
		default:
		}

		respMsg, _ := proto.Marshal(&wsPB.Response{Payload: &wsPB.Response_Response{Response: &resp}})
		_ = conn.WriteMessage(websocket.BinaryMessage, respMsg)
	})
	return upgrader
}

func (ws *WebsocketServer) handleWebsocketConnect(conn *websocket.Conn) {
	ws.logger.Debug().
		Any("Client connected", conn.RemoteAddr()).
		Send()

	// NOTE:
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 60))
}

func (ws *WebsocketServer) handleWebsocketDisconnect(conn *websocket.Conn, err error) {
	// TODO: why "EOF" error?
	if err != nil {
		ws.logger.Debug().AnErr("Failed to close connection", err).Send()
	}
	clientConn := conn.Session().(*ConnectionData)
	clientConn.ConnectionPool.activeConnections.Delete(clientConn.ID)

	ws.logger.Debug().Any("client disconnected", conn.RemoteAddr()).Send()

	// 	connData := conn.Session().(*ConnectionData)
	// 	connData.mu.Lock()
	// 	ws.msgScheduler.wheel.Remove()
	// 	for subscription := range connData.subscriptions {
	// 		connData.
	// 	}
	// 	connData.mu.Unlock()
}

func (ws *WebsocketServer) handlePublish(channelIdx wsPB.Channel, msg []byte, conn *websocket.Conn) error {
	err := ws.validatePubSubChannel(channelIdx)
	if err != nil {
		return err
	}
	connData := conn.Session().(*ConnectionData)
	channel := connData.ConnectionPool.channels[channelIdx]
	if channel == nil {
		return fmt.Errorf("not allowed channel: %s", channelIdx)
	}
	err = channel.Publish(connData.ID, msg)
	if err != nil {
		return fmt.Errorf("failed to publish message")
	}
	return nil
}

func (ws *WebsocketServer) handleSubscribe(channelIdx wsPB.Channel, publisherIDs []uint32, conn *websocket.Conn) ([]byte, error) {
	err := ws.validatePubSubChannel(channelIdx)
	if err != nil {
		return nil, err
	}

	connData := conn.Session().(*ConnectionData)
	connData.mu.Lock()
	subChannel := connData.ConnectionPool.channels[channelIdx]
	connData.subscriptions[channelIdx] = append(connData.subscriptions[channelIdx], publisherIDs...)
	connData.mu.Unlock()

	fetchedMessages, err := subChannel.GetMessages(publisherIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages")
	}

	ws.msgScheduler.scheduleConnection(
		ConnSlot{
			ConnData:       connData,
			Conn:           conn,
			ChannelToFlush: channelIdx,
		},
		ws.msgScheduler.calculateNextSchedule(),
	)

	return fetchedMessages, nil
}

func (ws *WebsocketServer) handleUnsubscribe(channelIdx wsPB.Channel, publisherIDs []uint32, conn *websocket.Conn) {
	err := ws.validatePubSubChannel(channelIdx)
	if err != nil {
		return
	}

	// Remove channel from subscription
	connData := conn.Session().(*ConnectionData)
	if len(publisherIDs) == 0 {
		connData.mu.Lock()
		// delete(connData.subscriptions, channelIdx)
		connData.subscriptions[channelIdx] = connData.subscriptions[channelIdx][:0]
		connData.mu.Unlock()
		return
	}

	// Swap all subscriptions in channel
	connData.mu.Lock()
	connData.subscriptions[channelIdx] = publisherIDs
	connData.mu.Unlock()
}

func (ws *WebsocketServer) validatePubSubChannel(channelIdx wsPB.Channel) error {
	_, ok := wsPB.Channel_name[int32(channelIdx)]
	if !ok {
		return fmt.Errorf("no such channel: %s", channelIdx)
	}

	return nil
}

func sendErrorResponse(conn *websocket.Conn, requestID string, code uint32, errMsg string) {
	resp := wsPB.ResponseMessage{
		RequestId:    requestID,
		StatusCode:   code,
		ErrorMessage: errMsg,
	}
	respMsg, _ := proto.Marshal(&wsPB.Response{Payload: &wsPB.Response_Response{Response: &resp}})
	_ = conn.WriteMessage(websocket.BinaryMessage, respMsg)
}
