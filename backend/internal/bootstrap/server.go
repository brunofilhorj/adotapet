package bootstrap

import (
	"log/slog"
	"net/http"
	"time"

	httpadapter "adotapet/internal/adapters/inbound/http"
)

func newServer(cfg Config, log *slog.Logger, services ...Service) *http.Server {
	return &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           httpadapter.NewRouter(cfg, log, services...),
		ReadHeaderTimeout: 5 * time.Second,
	}
}
