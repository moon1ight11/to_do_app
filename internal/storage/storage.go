package storage

import (
	"database/sql"
	"fmt"
	"todoapp/internal/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

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
		return nil, fmt.Errorf("storage.NewStorage: open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("storage.NewStorage: ping db: %w", err)
	}

	return &DataBase{
		DB:            db,
		MigrationsDir: cfg.Database.MigrationsDir,
	}, nil
}

func (d *DataBase) UpMigrations() error {
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("storage.UpMigrations: set dialect: %w", err)
	}

	if err := goose.Up(d.DB, d.MigrationsDir); err != nil {
		return fmt.Errorf("storage.UpMigrations: run migrations: %w", err)
	}

	return nil
}
