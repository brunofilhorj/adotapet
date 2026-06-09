package favorites

import (
	"net/http"

	httperrors "adotapet/internal/adapters/inbound/http/errors"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/puppies/{id}/favorite", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Adicionar favorito")
	})
	mux.HandleFunc("DELETE /api/v1/puppies/{id}/favorite", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Remover favorito")
	})
	mux.HandleFunc("GET /api/v1/me/favorites", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Listagem de favoritos")
	})
	mux.HandleFunc("POST /api/v1/me/saved-searches", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Criacao de busca salva")
	})
	mux.HandleFunc("GET /api/v1/me/saved-searches", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Listagem de buscas salvas")
	})
	mux.HandleFunc("DELETE /api/v1/me/saved-searches/{id}", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Remocao de busca salva")
	})
}
