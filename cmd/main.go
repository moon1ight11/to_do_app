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
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("main: panic recovered: %v\n", r)
			os.Exit(1)
		}
	}()

	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("main: failed to load config: %v", err))
	}

	deps := app.InitDependencies(cfg)
	defer deps.Close()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: deps.Router.GetEngine(),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	serverError := make(chan error, 1)

	go func() {
		deps.Logger.Info("main: server starting", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			deps.Logger.Error("main: server error", "error", err)
			serverError <- err
		}
	}()

	select {
	case <-quit:
		deps.Logger.Info("main: shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			deps.Logger.Error("main: forced shutdown", "error", err)
		}
	case err := <-serverError:
		deps.Logger.Fatal("main: server fatal error", "error", err)
	}

	deps.Logger.Info("main: server exited")
}
