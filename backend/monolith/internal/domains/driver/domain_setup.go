package driver

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"monolith/internal/domains/driver/application/command"
	"monolith/internal/domains/driver/application/query"
	"monolith/internal/domains/driver/infrastructure/postgres"
	"monolith/internal/domains/driver/infrastructure/postgres/sqlc"
	driverHTTP "monolith/internal/domains/driver/interface/http"
)

type DriverDomain struct {
	Commands DriverCommands
	Queries  DriverQueries
}

type DriverCommands struct {
	CreateDriver *command.CreateDriverHandler
}

type DriverQueries struct {
	GetDriverByUserID  *query.GetDriverByUserIDHandler
	GetFilteredDrivers *query.GetFilteredDriversHandler
}

func NewDriverDomain(db *pgxpool.Pool) *DriverDomain {
	driverRepo := postgres.NewDriverRepository(sqlc.New(db))

	createHandler := command.NewCreateDriverHandler(driverRepo)
	getDriverByUserIDHandler := query.NewGetDriverByUserIDHandler(driverRepo)
	getFilteredDriversHandler := query.NewGetFilteredDriversHandler(driverRepo)

	driverDomain := &DriverDomain{
		Commands: DriverCommands{
			CreateDriver: createHandler,
		},
		Queries: DriverQueries{
			GetDriverByUserID:  getDriverByUserIDHandler,
			GetFilteredDrivers: getFilteredDriversHandler,
		},
	}
	return driverDomain
}

func (d *DriverDomain) RegisterHTTPRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	driverHandler := driverHTTP.NewDriverHandler(
		d.Commands.CreateDriver,
		d.Queries.GetDriverByUserID,
		d.Queries.GetFilteredDrivers,
	)
	driverHTTP.RegisterDriverRoutes(router, driverHandler, authMiddleware)
}
