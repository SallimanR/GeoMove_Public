package server

import (
	"context"
	"sync"
	"time"
)

type ConnectionState string

const (
	StateDisconnected ConnectionState = "disconnected"
	StateConnecting   ConnectionState = "connecting"
	StateConnected    ConnectionState = "connected"
	StateDegraded     ConnectionState = "degraded"
)

type Connection interface {
	Connect(ctx context.Context) error
	GetConnection(ctx context.Context) error
	Ping(ctx context.Context) error
	Disconnect(ctx context.Context) error
}

type ConnectionService struct {
	conn      Connection
	connError error
	stopCh    chan struct{}

	connectionTimeout     time.Duration
	connectionRetryPeriod time.Duration
	pingTimeout           time.Duration
	pingRetryPeriod       time.Duration

	mu sync.RWMutex
}

func NewConnectionService(
	conn Connection,
	connectionTimeout time.Duration,
	connectionRetryPeriod time.Duration,
	pingTimeout time.Duration,
	pingRetryPeriod time.Duration,
) *ConnectionService {
	return &ConnectionService{
		conn:                  conn,
		connectionTimeout:     connectionTimeout,
		connectionRetryPeriod: connectionRetryPeriod,
		pingTimeout:           pingTimeout,
		pingRetryPeriod:       pingRetryPeriod,
	}
}

func (cm *ConnectionService) Start(ctx context.Context) {
	cm.attemptConnection(ctx)

	go cm.monitorConnection(ctx)
}

func (cm *ConnectionService) attemptConnection(ctx context.Context) {
	connectCtx, cancel := context.WithTimeout(ctx, cm.connectionTimeout)
	err := cm.conn.Connect(connectCtx)
	cancel()

	cm.mu.Lock()
	cm.connError = err
	cm.mu.Unlock()
}

func (cm *ConnectionService) checkConnection(ctx context.Context) {
	pingCtx, cancel := context.WithTimeout(ctx, cm.pingTimeout)
	err := cm.conn.Ping(pingCtx)
	cancel()

	cm.mu.Lock()
	cm.connError = err
	cm.mu.Unlock()
}

func (cm *ConnectionService) monitorConnection(ctx context.Context) {
	ticker := time.NewTicker(cm.connectionRetryPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-cm.stopCh:
			return
		case <-ticker.C:
			cm.checkConnection(ctx)
		}
	}
}
