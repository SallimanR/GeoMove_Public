package setup

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"monolith/internal/domains/geolocation/application/command"
	"monolith/internal/domains/geolocation/application/query"
	"monolith/internal/domains/geolocation/infrastructure/postgres"
	"monolith/internal/domains/geolocation/infrastructure/postgres/sqlc"
	geoHTTP "monolith/internal/domains/geolocation/interface/http"
	"monolith/internal/domains/geolocation/interface/websocket"
	"monolith/internal/websockethub"
)

type GeolocationDomain struct {
	Commands GeolocationCommands
	Queries  GeolocationQueries
}

type GeolocationCommands struct {
	CreateDriverRealtime *command.CreateDriverRealtimeHandler
	UpdateDriverRealtime *command.UpdateDriverRealtimeHandler
}

type GeolocationQueries struct {
	GetDriverByID              query.GetDriverRealtimeByIDHandler
	FindClosestDriversRealtime query.FindClosestDriversRealtimeHandler
}

func NewGeolocationDomain(db *pgxpool.Pool, wsServer *websockethub.WebsocketServer, logger zerolog.Logger) (*GeolocationDomain, error) {
	geoRepo := postgres.NewGeolocationRepository(sqlc.New(db))

	createCmd := command.NewCreateDriverRealtimeHandler(geoRepo)
	updateCmd := command.NewUpdateDriverRealtimeHandler(geoRepo)

	getDriverQuery := query.NewGetDriverRealtimeByIDHandler(geoRepo)
	findClosestQuery := query.NewFindClosestDriversRealtimeHandler(geoRepo)

	_, _, err := websocket.NewGPSRealtimeChannel(wsServer, []string{"tow_driver", "tow_subscriber"}, updateCmd, logger)
	if err != nil {
		return nil, err
	}

	return &GeolocationDomain{
		Commands: GeolocationCommands{
			CreateDriverRealtime: createCmd,
			UpdateDriverRealtime: updateCmd,
		},
		Queries: GeolocationQueries{
			GetDriverByID:              getDriverQuery,
			FindClosestDriversRealtime: findClosestQuery,
		},
	}, nil
}

func (d *GeolocationDomain) RegisterHTTPRoutes(router *gin.RouterGroup) {
	geoHandler := geoHTTP.NewGeolocationHandler(d.Queries.FindClosestDriversRealtime)
	geoHTTP.RegisterGeolocationRoutes(router, geoHandler)
}
