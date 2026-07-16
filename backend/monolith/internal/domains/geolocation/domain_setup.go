package geolocation

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"monolith/internal/domains/geolocation/application/command"
	"monolith/internal/domains/geolocation/application/query"
	"monolith/internal/domains/geolocation/domain/entity"
	"monolith/internal/domains/geolocation/infrastructure/postgres"
	"monolith/internal/domains/geolocation/infrastructure/postgres/sqlc"
	geoHTTP "monolith/internal/domains/geolocation/interface/http"
	"monolith/internal/domains/geolocation/interface/websocket"
	"monolith/internal/websockethub"
)

const (
	staleDriverCutoff          = "2 minutes"
	staleDriverCleanupInterval = 30 * time.Second
)

type GeolocationDomain struct {
	Commands   GeolocationCommands
	Queries    GeolocationQueries
	gpsChannel *websockethub.PubSubChannel[entity.MovingDriverWithPoints]

	stopCleanup chan struct{}
}

type GeolocationCommands struct {
	UpdateMovingDriver       *command.UpdateMovingDriverHandler
	DeleteStaleMovingDrivers *command.DeleteStaleMovingDriversHandler
}

type GeolocationQueries struct {
	GetMovingDriversByID                query.GetMovingDriverByIDHandler
	GetClosestWithinRadiusMovingDrivers query.GetClosestWithinRadiusMovingDriversHandler
}

func NewGeolocationDomain(db *pgxpool.Pool, wsServer *websockethub.WebsocketServer, logger zerolog.Logger) (*GeolocationDomain, error) {
	geoRepo := postgres.NewGeolocationRepository(sqlc.New(db))

	updateCmd := command.NewUpdateMovingDriverHandler(geoRepo)
	deleteStaleCmd := command.NewDeleteStaleMovingDriversHandler(geoRepo)
	getDriverQuery := query.NewGetMovingDriverByIDHandler(geoRepo)
	getClosestQuery := query.NewGetClosestWithinRadiusMovingDriversHandler(geoRepo)

	gpsChan, _, err := websocket.NewGPSRealtimeChannel(wsServer, []string{"tow_driver", "tow_subscriber"}, updateCmd, logger)
	if err != nil {
		return nil, err
	}

	domain := &GeolocationDomain{
		Commands: GeolocationCommands{
			UpdateMovingDriver:       updateCmd,
			DeleteStaleMovingDrivers: deleteStaleCmd,
		},
		Queries: GeolocationQueries{
			GetMovingDriversByID:                getDriverQuery,
			GetClosestWithinRadiusMovingDrivers: getClosestQuery,
		},
		gpsChannel:  gpsChan.Channel,
		stopCleanup: make(chan struct{}),
	}

	go domain.runStaleDriverCleanup(logger)

	return domain, nil
}

func (d *GeolocationDomain) RegisterHTTPRoutes(router *gin.RouterGroup) {
	geoHandler := geoHTTP.NewGeolocationHandler(d.Queries.GetClosestWithinRadiusMovingDrivers, d.gpsChannel)
	geoHTTP.RegisterGeolocationRoutes(router, geoHandler)
}

func (d *GeolocationDomain) Stop() {
	close(d.stopCleanup)
}

func (d *GeolocationDomain) runStaleDriverCleanup(logger zerolog.Logger) {
	ticker := time.NewTicker(staleDriverCleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := d.Commands.DeleteStaleMovingDrivers.Handle(context.Background(), staleDriverCutoff)
			if err != nil {
				logger.Err(err).Msg("failed to cleanup stale moving drivers")
			}
		case <-d.stopCleanup:
			return
		}
	}
}
