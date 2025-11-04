//go:build wireinject
// +build wireinject

package wire

import (
	"hona/backend/bootstrap"

	"hona/backend/internal/application/service"
	"hona/backend/internal/application/usecase"
	domainjwt "hona/backend/internal/domain/jwt"
	"hona/backend/internal/domain/ports"
	"hona/backend/internal/infrastructure/jwt"
	"hona/backend/internal/infrastructure/persistence"
	"hona/backend/internal/infrastructure/persistence/seeder"
	"hona/backend/internal/infrastructure/storage"
	"hona/backend/internal/presentation/controllers/v1/admin"
	"hona/backend/internal/presentation/controllers/v1/general"
	"hona/backend/internal/presentation/controllers/v1/user"
	"hona/backend/internal/presentation/middleware"

	"github.com/google/wire"
)

var StorageProviderSet = wire.NewSet(
	storage.NewS3Storage,
	wire.Struct(new(Storage), "*"),
)

var RepositoryProviderSet = wire.NewSet(
	persistence.NewRepositoryFactory,
	persistence.NewUnitOfWork,
	persistence.NewPostgresDatabase,
	wire.Bind(new(ports.RepositoryFactory), new(*persistence.RepositoryFactory)),
	wire.Bind(new(ports.UnitOfWork), new(*persistence.UnitOfWork)),
)
var ServiceProviderSet = wire.NewSet(
	service.NewUserService,
	jwt.NewJWTService,
	jwt.NewJWTKeyManager,
	service.NewRBACService,
	service.NewPetService,
	wire.Bind(new(domainjwt.JWTService), new(*jwt.JWTService)),
	wire.Bind(new(domainjwt.JWTKeyManager), new(*jwt.JWTKeyManager)),
	wire.Bind(new(usecase.RBACService), new(*service.RBACService)),
	wire.Bind(new(usecase.UserService), new(*service.UserService)),
	wire.Bind(new(usecase.PetService), new(*service.PetService)),
)

var GeneralControllersProviderSet = wire.NewSet(
	general.NewGeneralUserController,
	wire.Struct(new(GeneralControllers), "*"),
)

var AdminControllersProviderSet = wire.NewSet(
	admin.NewAdminRBACController,
	wire.Struct(new(AdminControllers), "*"),
)

var UserControllersProviderSet = wire.NewSet(
	user.NewUserPetController,
	wire.Struct(new(UserControllers), "*"),
)

var ControllersProviderSet = wire.NewSet(
	wire.Struct(new(Controllers), "*"),
)

var MiddlewaresProviderSet = wire.NewSet(
	middleware.NewLocalizationMiddleware,
	middleware.NewRecoveryMiddleware,
	middleware.NewRBACMiddleware,
	middleware.NewAuthMiddleware,
	wire.Struct(new(Middlewares), "*"),
)

var SeederProviderSet = wire.NewSet(
	seeder.NewDatabaseSeeder,
	wire.Struct(new(Seeder), "*"),
)

var ProviderSet = wire.NewSet(
	MiddlewaresProviderSet,
	ControllersProviderSet,
	GeneralControllersProviderSet,
	AdminControllersProviderSet,
	UserControllersProviderSet,
	ServiceProviderSet,
	RepositoryProviderSet,
	SeederProviderSet,
	StorageProviderSet,
)

type GeneralControllers struct {
	GeneralUserController *general.GeneralUserController
}

type AdminControllers struct {
	AdminRBACController *admin.AdminRBACController
}

type UserControllers struct {
	UserPetController *user.UserPetController
}

type Controllers struct {
	GeneralControllers *GeneralControllers
	AdminControllers   *AdminControllers
	UserControllers    *UserControllers
}

type Middlewares struct {
	LocalizationMiddleware *middleware.LocalizationMiddleware
	RecoveryMiddleware     *middleware.RecoveryMiddleware
	AuthMiddleware         *middleware.AuthMiddleware
	RBACMiddleware         *middleware.RBACMiddleware
}

type Seeder struct {
	DatabaseSeeder *seeder.DatabaseSeeder
}

type Storage struct {
	S3Storage *storage.S3Storage
}

type Seeder struct {
	DatabaseSeeder *seeder.DatabaseSeeder
}

type Application struct {
	Controllers *Controllers
	Middlewares *Middlewares
	Seeder      *Seeder
	Storage     *Storage
}

func NewApplication(controllers *Controllers, middlewares *Middlewares, seeder *Seeder, storage *Storage) *Application {
	return &Application{
		Controllers: controllers,
		Middlewares: middlewares,
		Seeder:      seeder,
		Storage:     storage,
	}
}

func InitializeApplication(container *bootstrap.Config) (*Application, error) {
	wire.Build(
		ProviderSet,
		NewApplication,
	)
	return &Application{}, nil
}
