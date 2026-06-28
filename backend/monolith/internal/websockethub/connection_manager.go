package websockethub

import (
	"sync"
	"sync/atomic"

	datastructures "monolith/pkg/data_structures"
)

type ConnectionData struct {
	ID             uint32
	ConnectionPool *ConnectionsByRole

	mu sync.RWMutex
	// subscriptions map[string][]uint32
	subscriptions [][]uint32
}

type ConnectionsByRole struct {
	activeConnections *datastructures.SyncMap[uint32, *ConnectionData]
	// connections       *datastructures.SyncMap[*websocket.Conn, *ConnectionData]
	stats struct {
		totalConns atomic.Int64
	}

	// channels map[string]ChannelActions
	channels []ChannelActions
}
