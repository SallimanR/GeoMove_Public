package websockethub

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"

	wsPB "monolith/internal/websockethub/proto"
	datastructures "monolith/pkg/data_structures"
)

type AuthSession interface {
	GetUserID() int64
	GetRoles() []string
}

type WebsocketServerOptions struct {
	Roles          []string       `json:"roles"`
	AllowedOrigins []string       `json:"-"`
	Logger         zerolog.Logger `json:"-"`
	DebugMode      bool           `json:"-"`
	TraceMode      bool           `json:"-"`
	TraceVerbose   bool           `json:"-"`
}

type WebsocketServer struct {
	upgrader          *websocket.Upgrader
	msgScheduler      MessageScheduler
	ConnectionsByRole map[string]*ConnectionsByRole

	logger         zerolog.Logger
	allowedOrigins map[string]struct{}
	debugMode      bool
	traceMode      bool
}

func NewWebsocketServer(config WebsocketServerOptions) *WebsocketServer {
	ws := &WebsocketServer{}

	allowedOrigins := make(map[string]struct{}, len(config.AllowedOrigins))
	for _, origin := range config.AllowedOrigins {
		allowedOrigins[origin] = struct{}{}
	}
	ws.allowedOrigins = allowedOrigins

	ws.upgrader = ws.newUpgrader()
	ws.msgScheduler = *NewMessageScheduler(10*time.Millisecond, 500)

	connectionsByRole := make(map[string]*ConnectionsByRole)
	for _, role := range config.Roles {
		connectionsByRole[role] = &ConnectionsByRole{
			activeConnections: &datastructures.SyncMap[uint32, *ConnectionData]{},
			channels:          make([]ChannelActions, len(wsPB.Channel_name)),
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

func (ws *WebsocketServer) WebsocketUpgradeHandler(ctx *gin.Context) {
	userVal, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}
	user, ok := userVal.(AuthSession)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	connRole := ctx.Param("role")
	if !hasRole(user.GetRoles(), connRole) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "role not permitted"})
		return
	}

	connID := uint32(user.GetUserID())

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
		subscriptions:  make([][]uint32, len(wsPB.Channel_name)),
	}
	conn.SetSession(connData)

	connPool.activeConnections.Store(connID, connData)
}

func hasRole(roles []string, target string) bool {
	for _, r := range roles {
		if r == target {
			return true
		}
	}
	return false
}

func (ws *WebsocketServer) newUpgrader() *websocket.Upgrader {
	upgrader := websocket.NewUpgrader()
	upgrader.CheckOrigin = func(r *http.Request) bool {
		if len(ws.allowedOrigins) == 0 {
			return true
		}
		origin := r.Header.Get("Origin")
		_, ok := ws.allowedOrigins[origin]
		return ok
	}

	upgrader.OnOpen(ws.handleWebsocketConnect)
	upgrader.OnClose(ws.handleWebsocketDisconnect)

	// If we care about error of writing back to connection, than retry. Else handling write error is NO-OP
	upgrader.OnMessage(func(conn *websocket.Conn, messageType websocket.MessageType, data []byte) {
		ws.logger.Debug().
			Int("bytes received", len(data)).
			Any("from client", conn.RemoteAddr()).
			Send()

		msg := &wsPB.Request{}
		err := proto.Unmarshal(data, msg)
		if err != nil {
			sendErrorResponse(conn, "", http.StatusBadRequest, "incorrect message")
			return
		}
		if msg.RequestId == "" {
			return
		}
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
