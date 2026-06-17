package users

import (
	"net/http"

	httperrors "adotapet/internal/adapters/inbound/http/errors"
	"adotapet/internal/adapters/inbound/http/middleware"
	"adotapet/internal/adapters/inbound/http/webserver"
)

type Service struct {
	authenticate webserver.Middleware
}

func NewService(accessTokenVerifier middleware.AccessTokenVerifier) Service {
	return Service{
		authenticate: middleware.Authenticate(accessTokenVerifier),
	}
}

func (s Service) Routes() []webserver.Route {
	return []webserver.Route{
		webserver.HandleFunc("GET /api/v1/me", handleGetMe, s.authenticate),
		webserver.HandleFunc("PATCH /api/v1/me/profile", handleUpdateProfile, s.authenticate),
	}
}

func handleGetMe(w http.ResponseWriter, r *http.Request) {
	if _, ok := middleware.AuthenticatedUser(r.Context()); !ok {
		httperrors.WriteJSON(w, http.StatusUnauthorized, httperrors.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Usuario autenticado ausente.",
		})
		return
	}
	httperrors.NotImplemented(w, "Consulta de perfil")
}

func handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	if _, ok := middleware.AuthenticatedUser(r.Context()); !ok {
		httperrors.WriteJSON(w, http.StatusUnauthorized, httperrors.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Usuario autenticado ausente.",
		})
		return
	}
	httperrors.NotImplemented(w, "Atualizacao de perfil")
}
