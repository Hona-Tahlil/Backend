//go:build wireinject
// +build wireinject

package wire

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/service"
	"hona/backend/internal/application/usecase"
	domainjwt "hona/backend/internal/domain/jwt"
	domainmail "hona/backend/internal/domain/mail"
	"hona/backend/internal/domain/ports"
	domainredis "hona/backend/internal/domain/ports/redis"
	domainstorage "hona/backend/internal/domain/storage"
	"hona/backend/internal/infrastructure/communication/mail"
	"hona/backend/internal/infrastructure/jwt"
	"hona/backend/internal/infrastructure/persistence"
	"hona/backend/internal/infrastructure/persistence/repository/redis"
	"hona/backend/internal/infrastructure/seeder"
	"hona/backend/internal/infrastructure/storage"
	"hona/backend/internal/infrastructure/websocket"
	"hona/backend/internal/presentation/controllers/v1/admin"
	"hona/backend/internal/presentation/controllers/v1/general"
	petsitter "hona/backend/internal/presentation/controllers/v1/pet_sitter"
	"hona/backend/internal/presentation/controllers/v1/user"
	"hona/backend/internal/presentation/middleware"

	"github.com/google/wire"
)

var StorageProviderSet = wire.NewSet(
	storage.NewS3Storage,
	wire.Bind(new(domainstorage.Storage), new(*storage.S3Storage)),
	wire.Struct(new(Storage), "*"),
)

var RepositoryProviderSet = wire.NewSet(
	persistence.NewRepositoryFactory,
	persistence.NewUnitOfWork,
	persistence.NewPostgresDatabase,
	persistence.NewRedisDatabase,
	redis.NewUserCacheRepository,
	wire.Bind(new(persistence.Cache), new(*persistence.RedisDatabase)),
	wire.Bind(new(domainredis.UserCacheRepository), new(*redis.UserCacheRepository)),
	wire.Bind(new(ports.RepositoryFactory), new(*persistence.RepositoryFactory)),
	wire.Bind(new(ports.UnitOfWork), new(*persistence.UnitOfWork)),
)
var ServiceProviderSet = wire.NewSet(
	wire.Struct(new(service.RequestServiceDeps), "*"),
	service.NewUserService,
	jwt.NewJWTService,
	jwt.NewJWTKeyManager,
	mail.NewEmailService,
	service.NewRBACService,
	service.NewPetService,
	service.NewRequestService,
	service.NewAddressService,
	service.NewPetSitterService,
	service.NewCommentService,
	service.NewChatService,
	service.NewWalletService,
	wire.Bind(new(domainjwt.JWTService), new(*jwt.JWTService)),
	wire.Bind(new(domainjwt.JWTKeyManager), new(*jwt.JWTKeyManager)),
	wire.Bind(new(usecase.RBACService), new(*service.RBACService)),
	wire.Bind(new(usecase.UserService), new(*service.UserService)),
	wire.Bind(new(usecase.PetService), new(*service.PetService)),
	wire.Bind(new(usecase.RequestService), new(*service.RequestService)),
	wire.Bind(new(usecase.AddressService), new(*service.AddressService)),
	wire.Bind(new(usecase.PetSitterService), new(*service.PetSitterService)),
	wire.Bind(new(usecase.CommentService), new(*service.CommentService)),
	wire.Bind(new(usecase.ChatService), new(*service.ChatService)),
	wire.Bind(new(usecase.WalletService), new(*service.WalletService)),
	wire.Bind(new(domainmail.Mail), new(*mail.EmailService)),
)

var GeneralControllersProviderSet = wire.NewSet(
	general.NewGeneralUserController,
	general.NewGeneralPetController,
	general.NewGeneralProvinceController,
	general.NewGeneralSearchController,
	wire.Struct(new(GeneralControllers), "*"),
)

var AdminControllersProviderSet = wire.NewSet(
	admin.NewAdminRBACController,
	admin.NewAdminPetSitterController,
	wire.Struct(new(AdminControllers), "*"),
)

var UserControllersProviderSet = wire.NewSet(
	user.NewUserPetController,
	user.NewUserRequestController,
	user.NewUserCommentController,
	user.NewUserChatController,
	user.NewUserProfileController,
	user.NewUserWalletController,
	wire.Struct(new(UserControllers), "*"),
)

var PetSitterControllersProviderSet = wire.NewSet(
	petsitter.NewPetSitterRegisterController,
	petsitter.NewPetSitterRequestController,
	petsitter.NewPetSitterChatController,
	petsitter.NewPetSitterSkillsController,
	petsitter.NewPetSitterCalendarController,
	petsitter.NewPetSitterWalletController,
	wire.Struct(new(PetSitterControllers), "*"),
)

var ControllersProviderSet = wire.NewSet(
	wire.Struct(new(Controllers), "*"),
)

var MiddlewaresProviderSet = wire.NewSet(
	middleware.NewLocalizationMiddleware,
	middleware.NewRecoveryMiddleware,
	middleware.NewRBACMiddleware,
	middleware.NewAuthMiddleware,
	middleware.NewCORSMiddleware,
	middleware.NewWebsocketMiddleware,
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
	PetSitterControllersProviderSet,
	ServiceProviderSet,
	RepositoryProviderSet,
	SeederProviderSet,
	StorageProviderSet,
)

type GeneralControllers struct {
	GeneralUserController          *general.GeneralUserController
	GeneralPetController           *general.GeneralPetController
	GeneralProvinceController      *general.GeneralProvinceController
	GeneralSearchController        *general.GeneralSearchController
}

type AdminControllers struct {
	AdminRBACController      *admin.AdminRBACController
	AdminPetSitterController *admin.AdminPetSitterController
}

type UserControllers struct {
	UserPetController     *user.UserPetController
	UserRequestController *user.UserRequestController
	UserCommentController *user.UserCommentController
	UserChatController    *user.UserChatController
	UserProfileController *user.UserProfileController
	UserWalletController  *user.UserWalletController
}

type PetSitterControllers struct {
	PetSitterRegisterController *petsitter.PetSitterRegisterController
	PetSitterRequestController  *petsitter.PetSitterRequestController
	PetSitterChatController     *petsitter.PetSitterChatController
	PetSitterSkillsController   *petsitter.PetSitterSkillsController
	PetSitterCalendarController *petsitter.PetSitterCalendarController
	PetSitterWalletController   *petsitter.PetSitterWalletController
}

type Controllers struct {
	GeneralControllers   *GeneralControllers
	AdminControllers     *AdminControllers
	UserControllers      *UserControllers
	PetSitterControllers *PetSitterControllers
}

type Middlewares struct {
	LocalizationMiddleware *middleware.LocalizationMiddleware
	RecoveryMiddleware     *middleware.RecoveryMiddleware
	AuthMiddleware         *middleware.AuthMiddleware
	RBACMiddleware         *middleware.RBACMiddleware
	CORSMiddleware         *middleware.CORSMiddleware
	WebsocketMiddleware    *middleware.WebsocketMiddleware
}

type Seeder struct {
	DatabaseSeeder *seeder.DatabaseSeeder
}

type Storage struct {
	S3Storage *storage.S3Storage
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

func InitializeApplication(container *bootstrap.Config, hub *websocket.Hub) (*Application, error) {
	wire.Build(
		ProviderSet,
		NewApplication,
	)
	return &Application{}, nil
}
