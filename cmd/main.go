package main

import (
	"log"
	"todoapp/internal/api"
	"todoapp/internal/api/handlers"
	"todoapp/internal/config"
	"todoapp/internal/api/jwt"
	"todoapp/internal/services"
	"todoapp/internal/storage"
	"todoapp/internal/storage/repos/settings"
	"todoapp/internal/storage/repos/tasks"
	"todoapp/internal/storage/repos/users"
)

func main() {
	// инициализация конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	// соединение с БД
	db, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatal("Failed to connect DB")
	}

	// применение миграций
	err = db.UpMigrations()
	if err != nil {
		log.Fatal("Failed to upping migrations", err)
	}

	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)

	// инициализация зависимостей
	userRepo := users.NewUserRepo(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, jwtService)
	authHandler := handlers.NewAuthHandler(userService, jwtService)

	settingsRepo := settings.NewSettingsRepo(db)
	settingsService := services.NewSettingsService(settingsRepo)
	settingsHandler := handlers.NewSettingsHandler(settingsService, jwtService)

	tasksRepo := tasks.NewTasksRepo(db)
	tasksService := services.NewTasksService(tasksRepo)
	tasksHandler := handlers.NewTasksHandler(tasksService, jwtService)

	// инициализация роутера
	router := api.NewRouter(userHandler, settingsHandler, tasksHandler, authHandler)

	// инициализация роутов
	router.Init(jwtService)

	// запуск роутера
	err = router.Run()
	if err != nil {
		log.Fatal("Failed to run Gin router", err)
	}
}
