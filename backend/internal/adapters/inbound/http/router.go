package http

import (
	"log/slog"
	"net/http"
	"time"

	"adotapet/internal/adapters/inbound/http/auth"
	"adotapet/internal/adapters/inbound/http/conversations"
	httperrors "adotapet/internal/adapters/inbound/http/errors"
	"adotapet/internal/adapters/inbound/http/favorites"
	"adotapet/internal/adapters/inbound/http/media"
	"adotapet/internal/adapters/inbound/http/middleware"
	"adotapet/internal/adapters/inbound/http/puppies"
	"adotapet/internal/adapters/inbound/http/users"
	"adotapet/internal/adapters/inbound/ws/chat"
	inport "adotapet/internal/app/port/in"
	"adotapet/internal/config"
)

func NewRouter(
	cfg config.Config,
	log *slog.Logger,
	registerUsers inport.RegisterUserInputPort,
	loginUsers inport.LoginInputPort,
	refreshTokens inport.RefreshTokenInputPort,
	verifyAccounts inport.VerifyAccountInputPort,
	resendVerification inport.ResendVerificationInputPort,
	accessTokenVerifier middleware.AccessTokenVerifier,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) {
		httperrors.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		httperrors.WriteJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})

	auth.RegisterRoutes(mux, registerUsers, loginUsers, refreshTokens, verifyAccounts, resendVerification)
	users.RegisterRoutes(mux, middleware.Authenticate(accessTokenVerifier))
	media.RegisterRoutes(mux)
	puppies.RegisterRoutes(mux)
	favorites.RegisterRoutes(mux)
	conversations.RegisterRoutes(mux)
	mux.HandleFunc("GET /ws/chat", chat.Handler())

	handler := middleware.RequestID(mux)
	handler = middleware.Recover(log)(handler)
	handler = middleware.AccessLog(log)(handler)
	handler = middleware.Timeout(30 * time.Second)(handler)
	handler = middleware.CORS(cfg.AppEnv)(handler)

	return handler
}
