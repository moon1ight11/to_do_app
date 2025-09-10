package main

import "todoapp/internal/db"

func main() {
	db.UppingMigrations()
	// здесь же должно быть подключение к БД и роутер
}