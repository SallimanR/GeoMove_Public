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
	FindDriverByID *query.GetDriverByIDHandler
}

func NewDriverDomain(db *pgxpool.Pool) *DriverDomain {
	driverRepo := postgres.NewDriverRepository(sqlc.New(db))
	createHandler := command.NewCreateDriverHandler(driverRepo)
	findByIDHandler := query.NewGetDriverByIDHandler(driverRepo)

	driverDomain := &DriverDomain{
		Commands: DriverCommands{
			CreateDriver: createHandler,
		},
		Queries: DriverQueries{
			FindDriverByID: findByIDHandler,
		},
	}
	return driverDomain
}

func (d *DriverDomain) RegisterHTTPRoutes(router *gin.RouterGroup) {
	driverHandler := driverHTTP.NewDriverHandler(d.Commands.CreateDriver, d.Queries.FindDriverByID)
	driverHTTP.RegisterDriverRoutes(router, driverHandler)
}
