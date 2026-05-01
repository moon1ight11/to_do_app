package storage

import (
	"database/sql"
	"fmt"
	"todoapp/internal/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// соединение с DB и return экземпляра
func NewStorage(cfg *config.Config) (*DataBase, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	datebase := DataBase{
		DB:            db,
		MigrationsDir: cfg.Database.MigrationsDir,
	}

	return &datebase, nil
}

// метод DB для применения миграций
func (d *DataBase) UpMigrations() error {
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	err := goose.Up(d.DB, d.MigrationsDir)
	if err != nil {
		return err
	}

	return nil
}
