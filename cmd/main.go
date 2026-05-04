package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todoapp/internal/app"
	"todoapp/internal/config"
)

func main() {
	// обработка паники
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic recovered: %v\n", r)
			os.Exit(1)
		}
	}()

	// инициализируем конфигурации
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// инициализируем все зависимости
	deps := app.InitDependencies(cfg)
	defer deps.Close()

	// Создаем сервер
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: deps.Router.GetEngine(),
	}

	// создаем каналы для сигналов завершения и ошибки сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	serverError := make(chan error, 1)

	// запуск сервера
	go func() {
		deps.Logger.Info("Server is starting", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			deps.Logger.Error("Failed to run server:", "error", err)
			serverError <- err
		}
	}()

	// ждем сигналы
	select {
	case <-quit:
		deps.Logger.Info("Shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			deps.Logger.Error("Server forced to shutdown:", "error", err)
		}
	case err := <-serverError:
		deps.Logger.Fatal("Server error:", "error", err)
	}

	deps.Logger.Info("Server exited")
}
