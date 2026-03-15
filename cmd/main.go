package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todoapp/internal/api"
	"todoapp/internal/api/handlers/authhandlers"
	"todoapp/internal/api/handlers/settingshandlers"
	"todoapp/internal/api/handlers/taskshandlers"
	"todoapp/internal/api/handlers/usershandlers"
	"todoapp/internal/api/jwt"
	"todoapp/internal/config"
	"todoapp/internal/services/settingsservice"
	"todoapp/internal/services/tasksservice"
	"todoapp/internal/services/usersservice"
	"todoapp/internal/storage"
	"todoapp/internal/storage/repos/settingsrepos"
	"todoapp/internal/storage/repos/tasksrepos"
	"todoapp/internal/storage/repos/usersrepos"
	"todoapp/pkg/logger"
)

func main() {
	// создаем логгер с записью в файл
	logger, err := logger.New("logs/ToDoApp.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	// инициализируем конфигурации
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config:", "error", err)
	}

	// соединяемся с БД
	db, err := storage.NewStorage(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to db:", "error", err)
	}
	defer db.DB.Close()

	// применяем миграций
	err = db.UpMigrations()
	if err != nil {
		logger.Fatal("Failed to upping migrations:", "error", err)
	}

	// инициализируем jwt
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)

	// инициализируем зависимости
	userRepo := usersrepos.NewUserRepo(db)
	userService := usersservice.NewUserService(userRepo)
	userHandler := usershandlers.NewUserHandler(userService, logger)
	authHandler := authhandlers.NewAuthHandler(userService, jwtService, logger)

	settingsRepo := settingsrepos.NewSettingsRepo(db)
	settingsService := settingsservice.NewSettingsService(settingsRepo)
	settingsHandler := settingshandlers.NewSettingsHandler(settingsService, logger)

	tasksRepo := tasksrepos.NewTasksRepo(db)
	tasksService := tasksservice.NewTasksService(tasksRepo)
	tasksHandler := taskshandlers.NewTasksHandler(tasksService, logger)

	// инициализируем роутер
	router := api.NewRouter(userHandler, settingsHandler, tasksHandler, authHandler)

	// инициализируем роуты
	router.Init(jwtService, logger)

	// Создаем HTTP сервер
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.GetEngine(),
	}

	// создаем каналы для сигналов завершения и ошибки сервера
	quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    serverError := make(chan error, 1)

	// запуск сервера
	 go func() {
        logger.Info("Server is starting on port :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error("Failed to run server:", "error", err)
            serverError <- err
        }
    }()

	// ждем сигналы
	select {
    case <-quit:
		// если поступил сигнал завершения - делаем шатдаун с таймаутом
        logger.Info("Shutting down server...")

        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        if err := srv.Shutdown(ctx); err != nil {
            logger.Error("Server forced to shutdown:", "error", err)
        }
    case err := <-serverError:
		// еслм пришла ошибка от сервера - фаталим
        logger.Fatal("Server error:", "error", err)
    }

    logger.Info("Server exited")
}
