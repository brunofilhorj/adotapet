package auth

import (
	"net/http"

	httperrors "adotapet/internal/adapters/inbound/http/errors"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/register", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Cadastro de usuario")
	})
	mux.HandleFunc("POST /api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Login")
	})
	mux.HandleFunc("POST /api/v1/auth/verify", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Verificacao de conta")
	})
	mux.HandleFunc("POST /api/v1/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Refresh token")
	})
}
