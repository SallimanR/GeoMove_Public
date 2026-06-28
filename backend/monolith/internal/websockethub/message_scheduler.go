package websockethub

import (
	"log"
	"net/http"
	"time"

	"github.com/lesismal/nbio/nbhttp/websocket"
	"google.golang.org/protobuf/proto"

	wsPB "monolith/internal/websockethub/proto"
	"monolith/pkg/data_structures/timewheel"
)

type ConnSlot struct {
	Conn     *websocket.Conn
	ConnData *ConnectionData
	// ChannelToFlush string
	ChannelToFlush wsPB.Channel
}

type MessageScheduler struct {
	wheel       *timewheel.TimeWheel[ConnSlot]
	base        int64
	periodTicks int64

	stop chan struct{}
}

func NewMessageScheduler(tick time.Duration, wheelSize int) *MessageScheduler {
	log.Println("periodTicks: ", int64(wheelSize*int(tick.Milliseconds())))
	return &MessageScheduler{
		wheel:       timewheel.NewTimeWheel[ConnSlot](tick, wheelSize),
		base:        0,
		periodTicks: int64(wheelSize) - 1,
		stop:        make(chan struct{}),
	}
}

func (ms *MessageScheduler) Start() {
	go ms.run()
}

func (ms *MessageScheduler) Stop() {
	close(ms.stop)
}

func (ms *MessageScheduler) run() {
	ticker := time.NewTicker(ms.wheel.Tick)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			due := ms.wheel.Advance()
			dueLen := len(due)
			if dueLen > 0 {
				nextExp := ms.calculateNextSchedule()
				log.Println("nextExp:", nextExp)
				// nextExp := time.Now().Add(5 * time.Second).UnixNano()
				for i := 0; i < dueLen; i++ {
					ms.flushConnection(due[i])
					// s.worker(slotData)

					ms.scheduleConnection(due[i], nextExp)
				}
			}
		case <-ms.stop:
			return
		}
	}
}

func (ms *MessageScheduler) scheduleConnection(connSlot ConnSlot, expires int64) {
	newSlot := ConnSlot{
		Conn:           connSlot.Conn,
		ConnData:       connSlot.ConnData,
		ChannelToFlush: connSlot.ChannelToFlush,
	}
	wheelSlot := ms.wheel.Add(newSlot, expires)
	log.Println("wheelSlot:", wheelSlot)
}

func (ms *MessageScheduler) calculateNextSchedule() int64 {
	nextExp := ms.base + (int64(ms.wheel.Current)+ms.periodTicks)*int64(ms.wheel.Tick)
	return nextExp
}

func (ms *MessageScheduler) flushConnection(connSlot ConnSlot) {
	resp := wsPB.ResponseMessage{}
	connSlot.ConnData.mu.Lock()
	publisherIDs := connSlot.ConnData.subscriptions[connSlot.ChannelToFlush]
	// channel, _ := connSlot.ConnData.ConnectionPool.channels[connSlot.ChannelToFlush]
	channel := connSlot.ConnData.ConnectionPool.channels[connSlot.ChannelToFlush]
	connSlot.ConnData.mu.Unlock()
	messages, err := channel.GetMessages(publisherIDs)
	if err != nil {
		// TODO: it is broadcast message, not req/resp
		// resp.StatusCode = http.StatusInternalServerError
		// resp.ErrorMessage = err.Error()
		// respMsg, err := proto.Marshal(&wsPB.Response{Payload: &wsPB.Response_Response{Response: &resp}})
		// if err != nil {
		// 	return
		// }
		// _ = connSlot.Conn.WriteMessage(websocket.BinaryMessage, respMsg)
		return
	}
	resp.StatusCode = http.StatusOK
	resp.Data = messages
	respMsg, err := proto.Marshal(&wsPB.Response{Payload: &wsPB.Response_Response{Response: &resp}})
	if err != nil {
		return
	}
	_ = connSlot.Conn.WriteMessage(websocket.BinaryMessage, respMsg)
}
