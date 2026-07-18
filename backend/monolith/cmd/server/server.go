package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"monolith/internal/auth"
	"monolith/internal/auth/sqlc"
	"monolith/internal/database"
	"monolith/internal/domains/driver"
	"monolith/internal/domains/geolocation"
	"monolith/internal/domains/order"
	"monolith/internal/notification"
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

	withDriverDomain      bool
	driverDomain          *driver.DriverDomain
	withGeolocationDomain bool
	geolocationDomain     *geolocation.GeolocationDomain
	withOrderDomain       bool
	orderDomain           *order.OrderDomain

	withAuth       bool
	authService    *auth.Service
	authMiddleware gin.HandlerFunc

	withNotification         bool
	notificationService      *notification.Service
	notificationContactEmail string

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
		s.httpRouter.Use(cors.New(newCORSConfig()))
	}

	if s.httpAPI == nil {
		s.httpAPI = s.httpRouter.Group("/api/v1")
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
		return nil
	}
}

func WithDriverDomain() Option {
	return func(s *Server) error {
		s.withDriverDomain = true
		return nil
	}
}

func (s *Server) setupDriverDomain() error {
	if s.geolocationDomain != nil {
		err := s.setupGeolocationDomain()
		if err != nil {
			return err
		}
	}
	domain := driver.NewDriverDomain(s.db)
	domain.RegisterHTTPRoutes(s.httpAPI, s.authMiddleware, s.authService)
	s.driverDomain = domain
	return nil
}

func (s *Server) GetDriverDomain() *driver.DriverDomain {
	return s.driverDomain
}

func WithGeolocationDomain() Option {
	return func(s *Server) error {
		s.withGeolocationDomain = true
		return nil
	}
}

func (s *Server) setupGeolocationDomain() error {
	if s.geolocationDomain != nil {
		return nil
	}

	domain, err := geolocation.NewGeolocationDomain(s.db, s.wsServer, s.logger)
	if err != nil {
		return err
	}
	domain.RegisterHTTPRoutes(s.httpAPI)
	s.geolocationDomain = domain
	return nil
}

func (s *Server) GetGeolocationDomain() *geolocation.GeolocationDomain {
	return s.geolocationDomain
}

func WithOrderDomain() Option {
	return func(s *Server) error {
		s.withOrderDomain = true
		return nil
	}
}

func (s *Server) setupOrderDomain() error {
	if s.orderDomain != nil {
		return nil
	}
	domain := order.NewOrderDomain(s.db, s.notificationService)
	domain.RegisterHTTPRoutes(s.httpAPI, s.authMiddleware)
	s.orderDomain = domain
	return nil
}

func (s *Server) GetOrderDomain() *order.OrderDomain {
	return s.orderDomain
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

	s.registerWSRoutes()

	if s.withNotification {
		err := s.setupNotification()
		if err != nil {
			return err
		}
	}

	if s.withGeolocationDomain {
		err := s.setupGeolocationDomain()
		if err != nil {
			return err
		}
	}
	if s.withDriverDomain {
		err := s.setupDriverDomain()
		if err != nil {
			return err
		}
	}
	if s.withOrderDomain {
		err := s.setupOrderDomain()
		if err != nil {
			return err
		}
	}

	go s.startWSServer()

	return nil
}

// TODO: graceful shutdown and properly exit from program
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down HTTP server...")
	if s.geolocationDomain != nil {
		s.geolocationDomain.Stop()
	}
	if s.notificationService != nil {
		err := s.notificationService.Close()
		if err != nil {
			return err
		}
	}
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

func loadCORSOrigins() []string {
	raw := os.Getenv("ALLOWED_CORS_ORIGINS")
	if raw == "" {
		return nil
	}
	parts := make([]string, 0)
	for _, p := range strings.Split(raw, ",") {
		if t := strings.TrimSpace(p); t != "" {
			parts = append(parts, t)
		}
	}
	return parts
}

func newCORSConfig() cors.Config {
	origins := loadCORSOrigins()
	if len(origins) == 0 {
		return cors.Config{
			AllowOriginFunc:  func(origin string) bool { return true },
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Cookie"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}
	}
	return cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Cookie"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

func loadAllowedOrigins() []string {
	raw := os.Getenv("ALLOWED_WS_ORIGINS")
	if raw == "" {
		return nil
	}
	parts := make([]string, 0)
	for _, p := range strings.Split(raw, ",") {
		if t := strings.TrimSpace(p); t != "" {
			parts = append(parts, t)
		}
	}
	return parts
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

func (s *Server) registerWSRoutes() {
	if s.wsServer == nil {
		allowedOrigins := loadAllowedOrigins()
		wsOptions := websockethub.WebsocketServerOptions{
			Roles:          []string{"tow_driver", "tow_subscriber"},
			AllowedOrigins: allowedOrigins,
			Logger:         s.logger,
		}
		s.wsServer = websockethub.NewWebsocketServer(wsOptions)
	}
	wsGroup := s.httpRouter.Group("/ws")
	if s.authMiddleware != nil {
		wsGroup.Use(s.authMiddleware)
	}
	wsGroup.GET("/:role", s.wsServer.WebsocketUpgradeHandler)
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

func WithNotification() Option {
	return func(s *Server) error {
		s.withNotification = true
		return nil
	}
}

func WithNotificationContactEmail(email string) Option {
	return func(s *Server) error {
		s.notificationContactEmail = email
		return nil
	}
}

func (s *Server) setupNotification() error {
	contactEmail := s.notificationContactEmail
	if contactEmail == "" {
		contactEmail = os.Getenv("VAPID_CONTACT_EMAIL")
	}

	svc, err := notification.NewService(notification.Config{
		DB:           s.db,
		ContactEmail: contactEmail,
		Logger:       s.logger,
	})
	if err != nil {
		return fmt.Errorf("setup notification service: %w", err)
	}

	notification.RegisterNotificationRoutes(s.httpAPI, svc.Handler, s.authMiddleware)
	s.notificationService = svc
	return nil
}

func (s *Server) GetNotificationService() *notification.Service {
	return s.notificationService
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
