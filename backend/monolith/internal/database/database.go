package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const (
	connectionTimeout     = 30 * time.Second
	connectionRetryPeriod = 10 * time.Second
	pingTimeout           = 5 * time.Second
)

type DBManager struct {
	config *pgxpool.Config
	pool   *pgxpool.Pool
	dbURL  string
	poolMu sync.RWMutex
}

func NewDBManager() (*DBManager, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("no .env file found or error loading: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}
	log.Printf("DATABASE_URL: %s", dbURL)

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database url: %v", err)
	}

	return &DBManager{
		dbURL:  dbURL,
		config: config,
	}, nil
}

func (dbm *DBManager) GetConnection() *pgxpool.Pool {
	dbm.poolMu.Lock()
	defer dbm.poolMu.Unlock()

	return dbm.pool
}

func (dbm *DBManager) Connect(ctx context.Context) error {
	log.Println("Attempting database connection")

	connnectCtx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(connnectCtx, dbm.config)
	if err != nil {
		return fmt.Errorf("Failed to get pool: connect to database: %v", err)
	}

	err = pool.Ping(connnectCtx)
	if err != nil {
		pool.Close()
		pool = nil
		return fmt.Errorf("Failed to ping database: %v", err)
	}

	dbm.poolMu.Lock()
	if dbm.pool != nil {
		dbm.pool.Close()
	}
	dbm.pool = pool
	dbm.poolMu.Unlock()
	return nil
}

func (dbm *DBManager) Ping(ctx context.Context) (bool, *ConnectionError) {
	// FIXME:
	dbm.poolMu.Lock()
	poolCopied := dbm.pool
	dbm.poolMu.Unlock()

	if poolCopied == nil {
		return false, &ConnectionError{Msg: "Database not connected", Retryable: true}
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	err := poolCopied.Ping(timeoutCtx)
	if err != nil {
		return false, &ConnectionError{Err: err, Msg: "Failed to ping db", Retryable: true}
	}
	return true, nil
}

func (dbm *DBManager) Disconnect() {
	dbm.poolMu.Lock()
	if dbm.pool != nil {
		// dbm.pool.Close()
		dbm.pool.Reset()
		dbm.pool = nil
	}
	dbm.poolMu.Unlock()
}

type ConnectionError struct {
	Err       error
	Msg       string
	Retryable bool
}
