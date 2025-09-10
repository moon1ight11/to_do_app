package db

import (
	"database/sql"
	"log"
	"github.com/pressly/goose/v3"
)

const migrationsDir = "./migrations"

var DB *sql.DB

func UppingMigrations() {
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Failed to set dialect: %v", err)
	}

	err := goose.Up(DB, migrationsDir)
	if err != nil {
		log.Fatalf("Failed to upping migrations: %v", err)
	}
}