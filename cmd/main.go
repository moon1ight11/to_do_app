package main

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
	"todoapp/internal/storage/repos/settingsrepos"
	"todoapp/internal/storage/repos/tasksrepos"
	"todoapp/internal/storage/repos/usersrepos"
	"todoapp/pkg/logger"
)

func main() {
	// создание логгера с записью в файл
	logger, err := logger.New("logs/ToDoApp.log")
	if err != nil {
		panic(err)
	}

	defer logger.Close()
	// инициализация конфигурации
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config:", err)
	}

	// соединение с БД
	db, err := storage.NewStorage(cfg)
	if err != nil {
		logger.Fatal("Failed to load config:", err)
	}

	// применение миграций
	err = db.UpMigrations()
	if err != nil {
		logger.Fatal("Failed to upping migrations:", err)
	}

	// инициализация jwt
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)

	// инициализация зависимостей
	userRepo := usersrepos.NewUserRepo(db)
	userService := usersservice.NewUserService(userRepo)
	userHandler := usershandlers.NewUserHandler(userService, logger)
	authHandler := authhandlers.NewAuthHandler(userService, jwtService, logger)

	settingsRepo := settingsrepos.NewSettingsRepo(db)
	settingsService := settingsservice.NewSettingsService(settingsRepo)
	settingsHandler := settingshandlers.NewSettingsHandler(settingsService, logger)

	tasksRepo := tasksrepos.NewTasksRepo(db)
	tasksService := tasksservice.NewTasksService(tasksRepo)
	tasksHandler := taskshandlers.NewTasksHandler(tasksService, logger)

	// инициализация роутера
	router := api.NewRouter(userHandler, settingsHandler, tasksHandler, authHandler)

	// инициализация роутов
	router.Init(jwtService, logger)

	// запуск роутера
	err = router.Run()
	if err != nil {
		logger.Fatal("Failed to run Gin router:", err)
	}
}
