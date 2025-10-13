package main

import (
	"log"
	"todoapp/internal/api"
	"todoapp/internal/api/handlers"
	"todoapp/internal/config"
	"todoapp/internal/services"
	"todoapp/internal/storage"
	"todoapp/internal/storage/repos/users"
	"todoapp/internal/storage/repos/settings"
	"todoapp/internal/storage/repos/tasks"
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
	
	// инициализация зависимостей
	userRepo := users.NewUserBase(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	settingsRepo := settings.NewSettingsBase(db)
	settingsService := services.NewSettingsService(settingsRepo)
	settingsHandler := handlers.NewSettingsHandler(settingsService)

	tasksRepo := tasks.NewTasksBase(db)
	tasksService := services.NewTasksService(tasksRepo)
	tasksHandler := handlers.NewTasksHandler(tasksService)	

	// запуск роутера
	api.UpRouter(userHandler, settingsHandler, tasksHandler)
}
