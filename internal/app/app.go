package app

import (
	"todoapp/internal/api"
	"todoapp/internal/api/handlers/authhandlers"
	"todoapp/internal/api/handlers/settingshandlers"
	"todoapp/internal/api/handlers/taskshandlers"
	"todoapp/internal/api/handlers/usershandlers"
	"todoapp/internal/api/jwt"
	"todoapp/internal/config"
	"todoapp/internal/services/settingsservice"
	"todoapp/internal/services/tasksservice"
	"todoapp/internal/services/usersservice"
	"todoapp/internal/storage"
	"todoapp/internal/storage/cache"
	"todoapp/internal/storage/repos/settingsrepos"
	"todoapp/internal/storage/repos/tasksrepos"
	"todoapp/internal/storage/repos/usersrepos"
	"todoapp/pkg/logger"
)

type Dependencies struct {
	Router *api.Router
	DB     *storage.DataBase
	Redis  *storage.RedisClient
	Logger logger.LoggerInterface
}

func InitDependencies(cfg *config.Config) *Dependencies {
	// создаем логгер
	log, err := logger.New(cfg)
	if err != nil {
		log.Fatal("Failed to init logger:", "error", err)
	}

	// соединяемся с БД
	db, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatal("Failed to connect to db:", "error", err)
	}

	// применяем миграции
	if err := db.UpMigrations(); err != nil {
		log.Fatal("Failed to upping migrations:", "error", err)
	}

	// подключаемся к редис
	redisClient, err := storage.NewRedisClient(cfg)
	if err != nil {
		log.Error("Failed to connect to Redis:", "error", err)
	}

	// инициализируем кэш
	var cacheService *cache.CacheService
	if redisClient != nil {
		cacheService = cache.NewCacheService(redisClient.Client)
	}

	// инициализируем jwt
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)

	// инициализируем зависимости
	userRepo := usersrepos.NewUserRepo(db)
	userService := usersservice.NewUserService(userRepo)
	userHandler := usershandlers.NewUserHandler(userService, log, cacheService)
	authHandler := authhandlers.NewAuthHandler(userService, jwtService, log)

	settingsRepo := settingsrepos.NewSettingsRepo(db)
	settingsService := settingsservice.NewSettingsService(settingsRepo)
	settingsHandler := settingshandlers.NewSettingsHandler(settingsService, log, cacheService)

	tasksRepo := tasksrepos.NewTasksRepo(db)
	tasksService := tasksservice.NewTasksService(tasksRepo)
	tasksHandler := taskshandlers.NewTasksHandler(tasksService, log, cacheService)

	// инициализируем роутер
	router := api.NewRouter(userHandler, settingsHandler, tasksHandler, authHandler)
	router.Init(jwtService, log)

	return &Dependencies{
		Router: router,
		DB:     db,
		Redis:  redisClient,
		Logger: log,
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
