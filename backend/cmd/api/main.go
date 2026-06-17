package main

import (
	"context"
	"log/slog"
	"os"

	"adotapet/internal/bootstrap"
	"adotapet/internal/config"
	"adotapet/internal/platform/logger"
)

func main() {
	os.Exit(run())
}

func run() int {
	cfg := config.Load()
	log := logger.New(cfg.AppEnv)

	app, err := bootstrap.NewApp(context.Background(), cfg, log)
	if err != nil {
		log.Error("failed to bootstrap app", slog.Any("error", err))
		return 1
	}
	defer func() {
		if err := app.Shutdown(context.Background()); err != nil {
			log.Error("failed to shutdown app resources", slog.Any("error", err))
		}
	}()

	log.Info("server running", slog.String("addr", app.Addr()), slog.String("env", cfg.AppEnv))

	if err := app.Run(context.Background()); err != nil {
		log.Error("server failed", slog.Any("error", err))
		return 1
	}

	return 0
}
