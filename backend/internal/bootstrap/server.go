package bootstrap

import (
	"log/slog"
	"net/http"
	"time"

	httpadapter "adotapet/internal/adapters/inbound/http"
)

func newServer(cfg Config, log *slog.Logger, authServices authServices) *http.Server {
	return &http.Server{
		Addr: cfg.HTTPAddr,
		Handler: httpadapter.NewRouter(
			cfg,
			log,
			authServices.registerUsers,
			authServices.loginUsers,
			authServices.refreshUserTokens,
			authServices.verifyAccounts,
			authServices.resendVerification,
			authServices.accessTokens,
		),
		ReadHeaderTimeout: 5 * time.Second,
	}
}
