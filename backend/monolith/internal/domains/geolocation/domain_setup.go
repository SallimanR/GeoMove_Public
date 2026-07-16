package geolocation

import (
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

type GeolocationDomain struct {
	Commands   GeolocationCommands
	Queries    GeolocationQueries
	gpsChannel *websockethub.PubSubChannel[entity.MovingDriverWithPoints]
}

type GeolocationCommands struct {
	UpdateMovingDriver *command.UpdateMovingDriverHandler
}

type GeolocationQueries struct {
	GetMovingDriversByID                query.GetMovingDriverByIDHandler
	GetClosestWithinRadiusMovingDrivers query.GetClosestWithinRadiusMovingDriversHandler
}

func NewGeolocationDomain(db *pgxpool.Pool, wsServer *websockethub.WebsocketServer, logger zerolog.Logger) (*GeolocationDomain, error) {
	geoRepo := postgres.NewGeolocationRepository(sqlc.New(db))

	updateCmd := command.NewUpdateMovingDriverHandler(geoRepo)
	getDriverQuery := query.NewGetMovingDriverByIDHandler(geoRepo)
	getClosestQuery := query.NewGetClosestWithinRadiusMovingDriversHandler(geoRepo)

	gpsChan, _, err := websocket.NewGPSRealtimeChannel(wsServer, []string{"tow_driver", "tow_subscriber"}, updateCmd, logger)
	if err != nil {
		return nil, err
	}

	return &GeolocationDomain{
		Commands: GeolocationCommands{
			UpdateMovingDriver: updateCmd,
		},
		Queries: GeolocationQueries{
			GetMovingDriversByID:                getDriverQuery,
			GetClosestWithinRadiusMovingDrivers: getClosestQuery,
		},
		gpsChannel: gpsChan.Channel,
	}, nil
}

func (d *GeolocationDomain) RegisterHTTPRoutes(router *gin.RouterGroup) {
	geoHandler := geoHTTP.NewGeolocationHandler(d.Queries.GetClosestWithinRadiusMovingDrivers, d.gpsChannel)
	geoHTTP.RegisterGeolocationRoutes(router, geoHandler)
}
