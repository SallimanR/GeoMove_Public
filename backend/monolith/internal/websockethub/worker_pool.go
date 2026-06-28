package websockethub

import "sync"

// TODO: Worker pool adaptive to number of OS threads
type WorkerPool struct {
	// batch       chan MessageBinary
	subscribers []uint32
	wg          sync.WaitGroup
}

type WorkerQueue[T any] struct {
	mu    sync.RWMutex
	items []T
}

// func (ws *WebsocketServer) checkConnections() {
// 	for i := range ws.roles {
// 		connPool := ws.connectionPools[ws.roles[i]].activeConnections
// 		for connectionIndx := range connPool {
// 		}
// 	}
// }
