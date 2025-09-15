package main

import (
	"log"
	"todoapp/internal/db"
)

func main() {
	// подключение к db
	datebase, err := db.DBConnection()
	if err != nil {
		log.Fatal("Failed to connecntion to DB:", err)
	}

	// создание экземпляра db
	db := db.NewDataBase(datebase, "./migrations")

	// применение миграций
	err = db.UppingMigrations()
	if err != nil {
		log.Fatal("Failed to upping migrations", err)
	}
}
