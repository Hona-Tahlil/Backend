//go:build wireinject
// +build wireinject

package wire

import (
	"hona/backend/bootstrap"

	"hona/backend/internal/application/service"
	"hona/backend/internal/infrastructure/jwt"
	"hona/backend/internal/infrastructure/persistence"
	"hona/backend/internal/presentation/controllers/v1/general"
	"hona/backend/internal/presentation/middleware"

	"github.com/google/wire"
)

var RepositoryProviderSet = wire.NewSet(
	persistence.NewRepositoryFactory,
	persistence.NewUnitOfWork,
	persistence.NewPostgresDatabase,
)

var ServiceProviderSet = wire.NewSet(
	service.NewUserService,
	jwt.NewJWTService,
	jwt.NewJWTKeyManager,
)

var GeneralControllersProviderSet = wire.NewSet(
	general.NewGeneralUserController,
	wire.Struct(new(GeneralControllers), "*"),
)

var ControllersProviderSet = wire.NewSet(
	wire.Struct(new(Controllers), "*"),
)

var MiddlewaresProviderSet = wire.NewSet(
	middleware.NewLocalizationMiddleware,
	middleware.NewRecoveryMiddleware,
	wire.Struct(new(Middlewares), "*"),
)

var ProviderSet = wire.NewSet(
	MiddlewaresProviderSet,
	ControllersProviderSet,
	GeneralControllersProviderSet,
	ServiceProviderSet,
	RepositoryProviderSet,
)

type GeneralControllers struct {
	GeneralUserController *general.GeneralUserController
}

type Controllers struct {
	GeneralControllers *GeneralControllers
}

type Middlewares struct {
	LocalizationMiddleware *middleware.LocalizationMiddleware
	RecoveryMiddleware     *middleware.RecoveryMiddleware
}

type Application struct {
	Controllers *Controllers
	Middlewares *Middlewares
}

func NewApplication(controllers *Controllers, middlewares *Middlewares) *Application {
	return &Application{
		Controllers: controllers,
		Middlewares: middlewares,
	}
}

func InitializeApplication(container *bootstrap.Config) (*Application, error) {
	wire.Build(
		ProviderSet,
		NewApplication,
	)
	return &Application{}, nil
}
