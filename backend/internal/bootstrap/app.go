package bootstrap

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// App owns the resources required to run the API.
type App struct {
	server    *http.Server
	resources []resource
	log       *slog.Logger
}

type resource interface {
	Shutdown(context.Context) error
}

// NewApp wires infrastructure, application services, and the HTTP server.
func NewApp(ctx context.Context, cfg Config, log *slog.Logger) (*App, error) {
	resources, err := prepareResources(ctx, cfg)
	if err != nil {
		return nil, err
	}

	authServices := newAuthServices(resources.database, cfg, log)

	return &App{
		server:    newServer(cfg, log, authServices),
		resources: resources.all(),
		log:       log,
	}, nil
}

// Addr returns the configured HTTP address.
func (app *App) Addr() string {
	return app.server.Addr
}

// Run starts the HTTP server and gracefully stops it on context cancellation or OS signals.
func (app *App) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		app.log.Info("starting api", slog.String("addr", app.server.Addr))
		if err := app.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := app.server.Shutdown(shutdownCtx); err != nil {
			return err
		}

		app.log.Info("api stopped")
		return nil
	}
}

// Shutdown closes resources initialized during bootstrap.
func (app *App) Shutdown(ctx context.Context) error {
	var errs []error

	for i := len(app.resources) - 1; i >= 0; i-- {
		if err := app.resources[i].Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
