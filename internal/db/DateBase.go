package db

import (
	"database/sql"
	"github.com/pressly/goose/v3"
	"log"
)

type App struct {
	DBex *DateBase
}

type DateBase struct {
	DB            *sql.DB
	MigrationsDir string
}

// экземпляр DB
func NewDataBase(db *sql.DB, migrationsDir string) *DateBase {
	datebase := DateBase{
		DB:            db,
		MigrationsDir: migrationsDir,
	}

	return &datebase
}

// метод DB для применения миграций
func (d *DateBase) UppingMigrations() error {
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Printf("Failed to set dialect: %v", err)
		return err
	}

	err := goose.Up(d.DB, d.MigrationsDir)
	if err != nil {
		log.Printf("Failed to upping migrations: %v", err)
		return err
	}

	log.Println("Migrations is upping")
	return nil
}

// соединение с DB
func DBConnection() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=fedor password=fedor_pass dbname=todo_app sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Failed to open database:", err)
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println("Failed to ping database:", err)
		return nil, err
	}

	log.Println("Successfully connected to DB")
	return db, nil
}
