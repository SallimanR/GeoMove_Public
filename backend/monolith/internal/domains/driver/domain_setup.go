package driver

import (
	"os"

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
	CreateDriver           *command.CreateDriverHandler
	UpdateProfileImage     *command.UpdateProfileImageHandler
	UpdateDriver           *command.UpdateDriverHandler
	CreateFreelyAvailable  *command.CreateFreelyAvailableHandler
	UpdateFreelyAvailable  *command.UpdateFreelyAvailableHandler
	DeleteFreelyAvailable  *command.DeleteFreelyAvailableHandler
}

type DriverQueries struct {
	GetDriverByUserID          *query.GetDriverByUserIDHandler
	GetFilteredDrivers         *query.GetFilteredDriversHandler
	GetFreelyAvailableByID     *query.GetFreelyAvailableByUserIDHandler
	GetFreelyAvailableDrivers  *query.GetFreelyAvailableDriversHandler
}

func NewDriverDomain(db *pgxpool.Pool) *DriverDomain {
	driverRepo := postgres.NewDriverRepository(sqlc.New(db))
	faRepo := postgres.NewFreelyAvailableRepository(sqlc.New(db))

	createHandler := command.NewCreateDriverHandler(driverRepo)
	updateProfileImageHandler := command.NewUpdateProfileImageHandler(driverRepo)
	updateDriverHandler := command.NewUpdateDriverHandler(driverRepo)
	getDriverByUserIDHandler := query.NewGetDriverByUserIDHandler(driverRepo)
	getFilteredDriversHandler := query.NewGetFilteredDriversHandler(driverRepo)

	createFAHandler := command.NewCreateFreelyAvailableHandler(faRepo)
	updateFAHandler := command.NewUpdateFreelyAvailableHandler(faRepo)
	deleteFAHandler := command.NewDeleteFreelyAvailableHandler(faRepo)
	getFAHandler := query.NewGetFreelyAvailableByUserIDHandler(faRepo)
	getFADriversHandler := query.NewGetFreelyAvailableDriversHandler(faRepo)

	driverDomain := &DriverDomain{
		Commands: DriverCommands{
			CreateDriver:          createHandler,
			UpdateProfileImage:    updateProfileImageHandler,
			UpdateDriver:          updateDriverHandler,
			CreateFreelyAvailable: createFAHandler,
			UpdateFreelyAvailable: updateFAHandler,
			DeleteFreelyAvailable: deleteFAHandler,
		},
		Queries: DriverQueries{
			GetDriverByUserID:         getDriverByUserIDHandler,
			GetFilteredDrivers:        getFilteredDriversHandler,
			GetFreelyAvailableByID:    getFAHandler,
			GetFreelyAvailableDrivers: getFADriversHandler,
		},
	}
	return driverDomain
}

func (d *DriverDomain) RegisterHTTPRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc, roleManager driverHTTP.UserRoleManager) {
	staticDir := os.Getenv("STATIC_DIR")

	driverHandler := driverHTTP.NewDriverHandler(
		d.Commands.CreateDriver,
		d.Commands.UpdateDriver,
		d.Queries.GetDriverByUserID,
		d.Queries.GetFilteredDrivers,
		d.Commands.UpdateProfileImage,
		d.Commands.CreateFreelyAvailable,
		d.Commands.UpdateFreelyAvailable,
		d.Commands.DeleteFreelyAvailable,
		d.Queries.GetFreelyAvailableByID,
		d.Queries.GetFreelyAvailableDrivers,
		staticDir,
		roleManager,
	)
	driverHTTP.RegisterDriverRoutes(router, driverHandler, authMiddleware)
}
