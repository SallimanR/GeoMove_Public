package setup

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"monolith/internal/domains/driver/application/command"
	"monolith/internal/domains/driver/application/query"
	"monolith/internal/domains/driver/infrastructure/db/postgres"
	"monolith/internal/domains/driver/infrastructure/db/sqlc"
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
	GetDriverByID      *query.GetDriverByIDHandler
	GetDriverByUserID  *query.GetDriverByUserIDHandler
	GetFilteredDrivers *query.GetFilteredDriversHandler
}

func NewDriverDomain(db *pgxpool.Pool) *DriverDomain {
	driverRepo := postgres.NewDriverRepository(sqlc.New(db))

	createHandler := command.NewCreateDriverHandler(driverRepo)
	getDriverByIDHandler := query.NewGetDriverByIDHandler(driverRepo)
	getDriverByUserIDHandler := query.NewGetDriverByUserIDHandler(driverRepo)
	getFilteredDriversHandler := query.NewGetFilteredDriversHandler(driverRepo)

	driverDomain := &DriverDomain{
		Commands: DriverCommands{
			CreateDriver: createHandler,
		},
		Queries: DriverQueries{
			GetDriverByID:      getDriverByIDHandler,
			GetDriverByUserID:  getDriverByUserIDHandler,
			GetFilteredDrivers: getFilteredDriversHandler,
		},
	}
	return driverDomain
}

func (d *DriverDomain) RegisterHTTPRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	driverHandler := driverHTTP.NewDriverHandler(
		d.Commands.CreateDriver,
		d.Queries.GetDriverByID,
		d.Queries.GetDriverByUserID,
		d.Queries.GetFilteredDrivers,
	)
	driverHTTP.RegisterDriverRoutes(router, driverHandler, authMiddleware)
}
