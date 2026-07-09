package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"monolith/internal/auth"
	"monolith/internal/auth/sqlc"
	"monolith/internal/database"
	"monolith/internal/domains/driver"
	"monolith/internal/websockethub"
)

const (
	connectionRetryPeriod = 5 * time.Second
	connectionTimeout     = 5 * time.Second

	defaultHTTPPort = 8100
)

// TODO: default options and default NewServer
type Server struct {
	httpAddrs  string
	httpRouter *gin.Engine
	httpAPI    *gin.RouterGroup

	httpServer *http.Server

	wsServer *websockethub.WebsocketServer

	// TODO: "with" for convenience or "without" for safety and explicitness
	// withoutDB bool
	withDB    bool
	dbManager *database.DBManager
	db        *pgxpool.Pool

	withDriverDomain bool
	driverDomain     *driver.DriverDomain

	withAuth       bool
	authService    *auth.Service
	authMiddleware gin.HandlerFunc

	logger zerolog.Logger
}

// NOTE: see: "Not using the functional options pattern (#11)" on https://100go.co/
type Option func(*Server) error

func NewServer(options ...Option) (*Server, error) {
	s := &Server{
		httpAddrs: fmt.Sprintf("127.0.0.1:%d", defaultHTTPPort),
	}

	if s.httpRouter == nil {
		s.httpRouter = gin.New()
		s.httpRouter.Use(gin.Logger())
		// if s.debugMode {
		// 	s.httpRouter.Use(gin.Logger())
		// }
	}

	// TODO:
	// s.logger =

	if s.httpAPI == nil {
		s.httpAPI = s.httpRouter.Group("/api/v1")
	}

	if s.wsServer == nil {
		wsOptions := websockethub.WebsocketServerOptions{
			Roles: []string{"tow_driver", "tow_subscriber"},
		}
		s.wsServer = websockethub.NewWebsocketServer(wsOptions)
		s.httpRouter.GET("/ws/:role/:id", s.wsServer.WebsocketUpgradeHandler)
	}

	// Apply options
	for _, option := range options {
		if err := option(s); err != nil {
			return nil, err
		}
	}

	if s.httpServer == nil {
		s.httpServer = &http.Server{
			Addr:    s.httpAddrs,
			Handler: s.httpRouter,
		}
	}

	if s.withDB && s.db == nil {
		var err error
		s.dbManager, err = database.NewDBManager()
		if err != nil {
			log.Fatal("Failed to setup DBManager")
		}
	}

	return s, nil
}

func WithHTTPPort(port int) Option {
	return func(s *Server) error {
		s.httpAddrs = fmt.Sprintf("127.0.0.1:%d", port)
		return nil
	}
}

func (s *Server) GetHttpAddrs() string {
	return s.httpAddrs
}

func WithDB() Option {
	return func(s *Server) error {
		s.withDB = true
		return nil
	}
}

func WithDBConnection(db *pgxpool.Pool) Option {
	return func(s *Server) error {
		s.db = db
		return nil
	}
}

func WithWebsocketServer(ws *websockethub.WebsocketServer) Option {
	return func(s *Server) error {
		s.wsServer = ws
		s.httpRouter.GET("/ws/:role/:id", s.wsServer.WebsocketUpgradeHandler)
		return nil
	}
}

func WithDriverDomain() Option {
	return func(s *Server) error {
		s.withDriverDomain = true
		return nil
	}
}

func (s *Server) setupDriverDomain() {
	domain := driver.NewDriverDomain(s.db)
	domain.RegisterHTTPRoutes(s.httpAPI, s.authMiddleware)
	s.driverDomain = domain
}

func (s *Server) GetDriverDomain() *driver.DriverDomain {
	return s.driverDomain
}

func WithLogger(logger zerolog.Logger) Option {
	return func(s *Server) error {
		s.logger = logger
		return nil
	}
}

func (s *Server) Start() error {
	s.setupMonitoringRoutes()
	go s.startListening()

	if s.db == nil && s.withDB {
		s.connectToDatabase()
	}

	if s.withAuth {
		s.setupAuth()
	}
	if s.withDriverDomain {
		s.setupDriverDomain()
	}

	go s.startWSServer()

	return nil
}

// TODO: graceful shutdown and properly exit from program
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down HTTP server...")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
		return err
	}
	log.Println("HTTP server shut down")
	return nil
}

func (s *Server) startWSServer() {
	s.wsServer.Run()
}

func (s *Server) setupMonitoringRoutes() {
	s.httpRouter.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})
	// TODO: statuses
	// s.httpRouter.GET("/:service_name", func(ctx *gin.Context) {
	// 	ctx.JSON(200, gin.H{
	// 		"status": s.serviceManager.services,
	// 	})
	// })
}

func (s *Server) setupAuth() {
	repo := auth.NewRepository(sqlc.New(s.db))

	cfg := auth.LoadConfig()
	log.Printf("Exchanging code with redirect_uri: %s", cfg.GoogleRedirectURL)
	providers := auth.GetOAuthProviders(cfg)

	service := auth.NewService(
		repo,
		providers,
		cfg.StaticDir,
	)

	s.authService = service
	s.authMiddleware = service.AuthMiddleware()

	auth.RegisterHTTPRoutes(s.httpAPI, service, s.authMiddleware)
}

func WithAuth() Option {
	return func(s *Server) error {
		s.withAuth = true
		return nil
	}
}

func (s *Server) startListening() {
	log.Printf("Server starting on: %s", s.httpAddrs)
	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed:", err)
	}
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return
	// 	}
	// }
}

func (s *Server) connectToDatabase() {
	for {
		connectCtx, cancel := context.WithTimeout(context.Background(), connectionRetryPeriod)
		err := s.dbManager.Connect(connectCtx)
		cancel()
		if err != nil {
			log.Printf("Failed connect to db: %s, reconnecting...", err)
			time.Sleep(connectionRetryPeriod)
			continue
		}
		s.db = s.dbManager.GetConnection()
		return
	}
}

// func (s *Server) monitorDatabaseConn() {
// 	for {
// 		connectCtx, cancel := context.WithTimeout(context.Background(), connectionRetryPeriod)
// 		ok, err := s.dbManager.Ping(connectCtx)
// 		if err != nil {}
// 		cancel()
// 	}
// }
