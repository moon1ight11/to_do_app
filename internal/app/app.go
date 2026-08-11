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
	l, err := logger.New(cfg)
	if err != nil {
		log.Fatalf("app.InitDependencies: logger.New: %v", err)
	}

	tracer, err := telemetry.Init(cfg)
	if err != nil {
		log.Fatalf("app.InitDependencies: telemetry.Init: %v", err)
	}

	m := metrics.NewMetrics()

	db, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatalf("app.InitDependencies: storage.NewStorage: %v", err)
	}

	if err := db.UpMigrations(); err != nil {
		log.Fatalf("app.InitDependencies: UpMigrations: %v", err)
	}

	redisClient, err := storage.NewRedisClient(cfg)
	if err != nil {
		log.Fatalf("app.InitDependencies: NewRedisClient: %v", err)
	}

	var cacheService *cache.CacheService
	if redisClient != nil {
		cacheService = cache.NewCacheService(redisClient.Client, tracer)
	}

	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)

	userRepo := usersrepos.NewUserRepo(db, tracer)
	userService := usersservice.NewUserService(userRepo, tracer)
	userHandler := usershandlers.NewUserHandler(userService, l, cacheService, tracer, m)
	authHandler := authhandlers.NewAuthHandler(userService, jwtService, l, tracer, m)

	settingsRepo := settingsrepos.NewSettingsRepo(db, tracer)
	settingsService := settingsservice.NewSettingsService(settingsRepo, tracer)
	settingsHandler := settingshandlers.NewSettingsHandler(settingsService, l, cacheService, tracer, m)

	tasksRepo := tasksrepos.NewTasksRepo(db, tracer)
	tasksService := tasksservice.NewTasksService(tasksRepo, tracer)
	tasksHandler := taskshandlers.NewTasksHandler(tasksService, l, cacheService, tracer, m)

	router := api.NewRouter(userHandler, settingsHandler, tasksHandler, authHandler)
	router.Init(jwtService, l, cfg, m)

	return &Dependencies{
		Router:    router,
		DB:        db,
		Redis:     redisClient,
		Logger:    l,
		Telemetry: tracer,
		Metrics:   m,
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
