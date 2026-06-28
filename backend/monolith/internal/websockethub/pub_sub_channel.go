package websockethub

import (
	"sync"
	"sync/atomic"

	"github.com/lesismal/nbio/nbhttp/websocket"
	"google.golang.org/protobuf/proto"

	wsPB "monolith/internal/websockethub/proto"
	datastructures "monolith/pkg/data_structures"
)

// TODO:
type MessageData struct {
	lastUpdated atomic.Int64
	ttl         atomic.Int64
	message     []byte
}

type ChannelActions interface {
	Publish(publisherID uint32, msg []byte) error
	GetMessages(publisherIDs []uint32) ([]byte, error)
}

type PubSubChannel struct {
	Messages *datastructures.SyncMap[uint32, []byte]

	muMessageQueue sync.Mutex
	messageQueue   []*[]byte

	muSubscribers sync.RWMutex
	subscribers   []*websocket.Conn // On delete don't shuffle but move last to deleted position
}

func NewPubSubChannel() *PubSubChannel {
	return &PubSubChannel{
		Messages: &datastructures.SyncMap[uint32, []byte]{},
	}
}

func (c *PubSubChannel) Publish() {}

func (c *PubSubChannel) GetMessages(publisherIDs []uint32) ([]byte, error) {
	messages := make([][]byte, 0, len(publisherIDs))
	for _, publisherID := range publisherIDs {
		message, ok := c.Messages.Load(publisherID)
		// TODO: handle !ok
		if !ok {
			continue
		}
		messages = append(messages, message)
	}

	var messagesProtobuf wsPB.MessageBatch
	messagesProtobuf.Data = messages
	return proto.Marshal(&messagesProtobuf)
}

func (mc *PubSubChannel) MessageBroadcast(msg []byte) {
	mc.muSubscribers.Lock()
	for i := range mc.subscribers {
		conn := mc.subscribers[i]
		_ = conn.WriteMessage(websocket.BinaryMessage, msg)
	}
	mc.muSubscribers.Unlock()
}
