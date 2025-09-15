package main

import (
	"log"
	"todoapp/internal/config"
	"todoapp/internal/db"
)

func main() {
	// инициализация конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	// подключение к db
	database, err := db.DBConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connecntion DB:", err)
	}

	// создание экземпляра db
	db := db.NewDataBase(database, cfg.Database.MigrationsDir)

	// применение миграций
	err = db.UppingMigrations()
	if err != nil {
		log.Fatal("Failed to upping migrations", err)
	}
}
