package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpadapter "adotapet/internal/adapters/inbound/http"
	"adotapet/internal/config"
	"adotapet/internal/platform/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.AppEnv)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           httpadapter.NewRouter(cfg, log),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info("starting api", slog.String("addr", cfg.HTTPAddr), slog.String("env", cfg.AppEnv))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("api stopped unexpectedly", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("api shutdown failed", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("api stopped")
}
