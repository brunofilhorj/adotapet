package http

import (
	"log/slog"
	"net/http"
	"time"

	httperrors "adotapet/internal/adapters/inbound/http/errors"
	"adotapet/internal/adapters/inbound/http/middleware"
	"adotapet/internal/adapters/inbound/http/webserver"
	"adotapet/internal/config"
)

func NewRouter(
	cfg config.Config,
	log *slog.Logger,
	services ...webserver.RouteProvider,
) http.Handler {
	return webserver.New(
		nil,
		webserver.WithRoutes(healthRoutes()...),
		webserver.WithServices(services...),
		webserver.WithMiddlewares(
			middleware.CORS(cfg.AppEnv),
			middleware.Timeout(30*time.Second),
			middleware.AccessLog(log),
			middleware.Recover(log),
			middleware.RequestID,
		),
	)
}

func healthRoutes() []webserver.Route {
	return []webserver.Route{
		webserver.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) {
			httperrors.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		}),
		webserver.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
			httperrors.WriteJSON(w, http.StatusOK, map[string]string{"status": "ready"})
		}),
	}
}
