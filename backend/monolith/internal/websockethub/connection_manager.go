package websockethub

import (
	"sync"
	"sync/atomic"

	datastructures "monolith/pkg/data_structures"
)

type ConnectionData struct {
	ID             int64
	ConnectionPool *ConnectionsByRole

	mu            sync.RWMutex
	subscriptions [][]int64
}

type ConnectionsByRole struct {
	activeConnections *datastructures.SyncMap[int64, *ConnectionData]
	stats             struct {
		totalConns atomic.Int64
	}

	channels []ChannelActions
}
