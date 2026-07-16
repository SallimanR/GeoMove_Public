package websockethub

import (
	"sync"

	"github.com/lesismal/nbio/nbhttp/websocket"

	datastructures "monolith/pkg/data_structures"
)

type ChannelActions interface {
	Publish(publisherID int64, msg []byte) error
	GetMessages(publisherIDs []int64) ([]byte, error)
}

type PubSubChannel[T any] struct {
	Messages *datastructures.SyncMap[int64, T]

	muMessageQueue sync.Mutex
	messageQueue   []*[]byte

	muSubscribers sync.RWMutex
	subscribers   []*websocket.Conn
}

func NewPubSubChannel[T any]() *PubSubChannel[T] {
	return &PubSubChannel[T]{
		Messages: &datastructures.SyncMap[int64, T]{},
	}
}

func (c *PubSubChannel[T]) GetMessages(publisherIDs []int64) []T {
	messages := make([]T, 0, len(publisherIDs))
	for _, publisherID := range publisherIDs {
		message, ok := c.Messages.Load(publisherID)
		if !ok {
			continue
		}
		messages = append(messages, message)
	}
	return messages
}

func (mc *PubSubChannel[T]) MessageBroadcast(msg []byte) {
	mc.muSubscribers.Lock()
	for i := range mc.subscribers {
		conn := mc.subscribers[i]
		_ = conn.WriteMessage(websocket.BinaryMessage, msg)
	}
	mc.muSubscribers.Unlock()
}
