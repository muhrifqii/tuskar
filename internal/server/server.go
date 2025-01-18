package server

import (
	"context"
	"os"

	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
	"github.com/muhrifqii/tuskar/internal/config"
	"github.com/muhrifqii/tuskar/internal/repository/postgresql"
	redisInternal "github.com/muhrifqii/tuskar/internal/repository/redis"
	"github.com/muhrifqii/tuskar/internal/rest"
	"github.com/muhrifqii/tuskar/internal/rest/middleware"
	"github.com/muhrifqii/tuskar/internal/rest/rest_utils"
	"github.com/muhrifqii/tuskar/usecase/authn"
	"github.com/muhrifqii/tuskar/usecase/provision"
	"github.com/muhrifqii/tuskar/usecase/task"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type (
	Server struct {
		app         *fiber.App
		args        ServerArgs
		provisioner *provision.Service
	}

	ServerArgs struct {
		Config      config.ApiConfig
		Logger      *zap.Logger
		Validator   *validator.Validate
		RedisClient *redis.Client
		DB          *sqlx.DB
	}
)

func NewServer(args ServerArgs) *Server {

	app := fiber.New(fiber.Config{
		CaseSensitive:            true,
		DisableHeaderNormalizing: true,
		JSONEncoder:              sonic.Marshal,
		JSONDecoder:              sonic.Unmarshal,
		ErrorHandler:             errorHandler,
	})

	// instantiate schema decoder
	encoderDecoder := rest_utils.SchemaEncoderDecoder{
		Encoder: schema.NewEncoder(),
		Decoder: schema.NewDecoder(),
	}

	// build redis client on fiber.Storage
	redisStorage := redisInternal.NewStorageRedis(args.RedisClient)

	// prepare middleware
	app.Use(middleware.Recover())
	app.Use(middleware.Cors(args.Config))
	app.Use(middleware.RequestID(args.Config))
	app.Use(middleware.Logger(args.Logger))
	app.Use(middleware.RateLimiter(50, redisStorage))
	app.Use(middleware.ActuatorHealthCheck())

	middleware.SetZapLogger(args.Logger)

	// prepare public route group
	apiPath := args.Config.ApiPrefix
	apiV1 := app.Group(apiPath)

	// prepare repository layer
	userRepository := postgresql.NewUserRepository(args.DB, args.Logger)
	taskRepository := postgresql.NewTaskRepository(args.DB, args.Logger)

	// build service layer
	authnSvc := authn.NewService(args.Logger, args.Config.JwtConfig, userRepository)
	provisioner := provision.NewService(args.Logger, userRepository)
	taskSvc := task.NewService(taskRepository, redisStorage, args.Logger)
	handlerParams := rest_utils.HandlerParams{
		Validator:            args.Validator,
		Logger:               args.Logger,
		SchemaEncoderDecoder: encoderDecoder,
	}
	// public handler
	rest.NewAuthnHandler(apiV1, authnSvc, handlerParams, args.Config.JwtConfig)

	// protected handler
	app.Use(middleware.RequireAuthn(args.Config.JwtConfig))
	rest.NewProtectedAuthnHandler(apiV1, authnSvc, handlerParams)
	rest.NewTaskHandler(apiV1, taskSvc, handlerParams)

	return &Server{
		app:         app,
		args:        args,
		provisioner: provisioner,
	}
}

func errorHandler(c *fiber.Ctx, err error) error {
	return rest_utils.ApiErrorResponseHandler(c, err)
}

func (s *Server) ProvisionSystemUser() {
	migrationUsername := os.Getenv("SYSTEM_USER_MIGRATION_USERNAME")
	migrationPassword := os.Getenv("SYSTEM_USER_MIGRATION_PASSWORD")

	err := s.provisioner.CreateSystemUser(context.Background(), migrationUsername, migrationPassword)
	if err != nil {
		s.args.Logger.Error("Failed to create system user", zap.Error(err))
	}
}

func (s *Server) Run() error {
	return s.app.Listen(s.args.Config.Port)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
