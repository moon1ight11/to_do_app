package app

import (
	"log"
	"todoapp/internal/api"
	"todoapp/internal/api/handlers/authhandlers"
	"todoapp/internal/api/handlers/settingshandlers"
	"todoapp/internal/api/handlers/taskshandlers"
	"todoapp/internal/api/handlers/usershandlers"
	"todoapp/internal/api/jwt"
	"todoapp/internal/config"
	"todoapp/internal/metrics"
	"todoapp/internal/services/settingsservice"
	"todoapp/internal/services/tasksservice"
	"todoapp/internal/services/usersservice"
	"todoapp/internal/storage"
	"todoapp/internal/storage/cache"
	"todoapp/internal/storage/repos/settingsrepos"
	"todoapp/internal/storage/repos/tasksrepos"
	"todoapp/internal/storage/repos/usersrepos"
	"todoapp/internal/telemetry"
	"todoapp/pkg/logger"

	"go.opentelemetry.io/otel/trace"
)

type Dependencies struct {
	Router    *api.Router
	DB        *storage.DataBase
	Redis     *storage.RedisClient
	Logger    logger.LoggerInterface
	Telemetry trace.Tracer
	Metrics   *metrics.Metrics
}

func InitDependencies(cfg *config.Config) *Dependencies {
	// создаем логгер
	logger, err := logger.New(cfg)
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	// инициализируем трейсер
	tracer, err := telemetry.Init(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}

	metrics := metrics.NewMetrics()

	// соединяемся с БД
	db, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	// применяем миграции
	if err := db.UpMigrations(); err != nil {
		log.Fatalf("Failed to upping migrations: %v", err)
	}

	// подключаемся к редис
	redisClient, err := storage.NewRedisClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// инициализируем кэш
	var cacheService *cache.CacheService
	if redisClient != nil {
		cacheService = cache.NewCacheService(redisClient.Client, tracer)
	}

	// инициализируем jwt
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)

	// инициализируем зависимости
	userRepo := usersrepos.NewUserRepo(db, tracer)
	userService := usersservice.NewUserService(userRepo, tracer)
	userHandler := usershandlers.NewUserHandler(userService, logger, cacheService, tracer, metrics)
	authHandler := authhandlers.NewAuthHandler(userService, jwtService, logger, tracer, metrics)

	settingsRepo := settingsrepos.NewSettingsRepo(db, tracer)
	settingsService := settingsservice.NewSettingsService(settingsRepo, tracer)
	settingsHandler := settingshandlers.NewSettingsHandler(settingsService, logger, cacheService, tracer, metrics)

	tasksRepo := tasksrepos.NewTasksRepo(db, tracer)
	tasksService := tasksservice.NewTasksService(tasksRepo, tracer)
	tasksHandler := taskshandlers.NewTasksHandler(tasksService, logger, cacheService, tracer, metrics)

	// инициализируем роутер
	router := api.NewRouter(userHandler, settingsHandler, tasksHandler, authHandler)
	router.Init(jwtService, logger, cfg, metrics)

	return &Dependencies{
		Router:    router,
		DB:        db,
		Redis:     redisClient,
		Logger:    logger,
		Telemetry: tracer,
		Metrics:   metrics,
	}
}

func (d *Dependencies) Close() {
	if d.DB != nil {
		d.DB.DB.Close()
	}
	if d.Redis != nil {
		d.Redis.Close()
	}
	if d.Logger != nil {
		d.Logger.Close()
	}
}
