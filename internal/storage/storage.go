package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"todoapp/internal/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
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
		log.Println("Failed to open database:", err)
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		log.Println("Failed to ping database:", err)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")

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

// соединение с редис
func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{Client: client}, nil
}

// метод для закрытия соединения с редис
func (r *RedisClient) Close() error {
    return r.Client.Close()
}