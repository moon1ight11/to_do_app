package db

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"log"
	"todoapp/internal/config"
)

type App struct {
	DBex *DataBase
}

type DataBase struct {
	DB            *sql.DB
	MigrationsDir string
}

// экземпляр DB
func NewDataBase(db *sql.DB, migrationsDir string) *DataBase {
	datebase := DataBase{
		DB:            db,
		MigrationsDir: migrationsDir,
	}

	return &datebase
}

// метод DB для применения миграций
func (d *DataBase) UppingMigrations() error {
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
func DBConnection(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s migratioinsDir=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.MigrationsDir,
	)

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
