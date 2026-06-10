package users

import (
	"net/http"

	httperrors "adotapet/internal/adapters/inbound/http/errors"
	"adotapet/internal/adapters/inbound/http/middleware"
)

func RegisterRoutes(mux *http.ServeMux, authenticate func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/me", authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := middleware.AuthenticatedUser(r.Context()); !ok {
			httperrors.WriteJSON(w, http.StatusUnauthorized, httperrors.ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "Usuario autenticado ausente.",
			})
			return
		}
		httperrors.NotImplemented(w, "Consulta de perfil")
	})))
	mux.Handle("PATCH /api/v1/me/profile", authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := middleware.AuthenticatedUser(r.Context()); !ok {
			httperrors.WriteJSON(w, http.StatusUnauthorized, httperrors.ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "Usuario autenticado ausente.",
			})
			return
		}
		httperrors.NotImplemented(w, "Atualizacao de perfil")
	})))
}
