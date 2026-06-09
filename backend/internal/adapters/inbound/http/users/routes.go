package users

import (
	"net/http"

	httperrors "adotapet/internal/adapters/inbound/http/errors"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/me", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Consulta de perfil")
	})
	mux.HandleFunc("PATCH /api/v1/me/profile", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Atualizacao de perfil")
	})
}
