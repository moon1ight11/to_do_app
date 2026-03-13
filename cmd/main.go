package main

import (
	"log"
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
)

func main() {
	// инициализация конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: %w", err)
	}

	// соединение с БД
	db, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatal("Failed to connect DB: %w", err)
	}

	// применение миграций
	err = db.UpMigrations()
	if err != nil {
		log.Fatal("Failed to upping migrations: %w", err)
	}

	// инициализация jwt
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)

	// инициализация зависимостей
	userRepo := usersrepos.NewUserRepo(db)
	userService := usersservice.NewUserService(userRepo)
	userHandler := usershandlers.NewUserHandler(userService)
	authHandler := authhandlers.NewAuthHandler(userService, jwtService)

	settingsRepo := settingsrepos.NewSettingsRepo(db)
	settingsService := settingsservice.NewSettingsService(settingsRepo)
	settingsHandler := settingshandlers.NewSettingsHandler(settingsService)

	tasksRepo := tasksrepos.NewTasksRepo(db)
	tasksService := tasksservice.NewTasksService(tasksRepo)
	tasksHandler := taskshandlers.NewTasksHandler(tasksService)

	// инициализация роутера
	router := api.NewRouter(userHandler, settingsHandler, tasksHandler, authHandler)

	// инициализация роутов
	router.Init(jwtService)

	// запуск роутера
	err = router.Run()
	if err != nil {
		log.Fatal("Failed to run Gin router: %w", err)
	}
}
