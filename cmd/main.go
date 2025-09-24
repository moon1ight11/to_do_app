package main

import "todoapp/internal/api"

// "log"
// "todoapp/internal/config"
// "todoapp/internal/storage"

func main() {
	// // инициализация конфигурации
	// cfg, err := config.Load()
	// if err != nil {
	// 	log.Fatal("Failed to load config", err)
	// }

	// // соединение с БД
	// db, err := storage.NewStorage(cfg)
	// if err != nil {
	// 	log.Fatal("Failed to connect DB")
	// }

	// // применение миграций
	// err = db.UpMigrations()
	// if err != nil {
	// 	log.Fatal("Failed to upping migrations", err)
	// }

	api.Router()
}
