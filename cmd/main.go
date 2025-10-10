package main

import (
	"log"
	"todoapp/internal/api"
	"todoapp/internal/api/handlers"
	"todoapp/internal/config"
	"todoapp/internal/services"
	"todoapp/internal/storage"
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
	
	// инициализация зависимостей
	userRepo := users.NewBase(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// запуск роутера
	api.UpRouter(userHandler)
}
