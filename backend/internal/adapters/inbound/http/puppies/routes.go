package puppies

import (
	"net/http"

	httperrors "adotapet/internal/adapters/inbound/http/errors"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/puppies", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Cadastro de filhote")
	})
	mux.HandleFunc("GET /api/v1/puppies/search", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Busca geolocalizada de filhotes")
	})
	mux.HandleFunc("GET /api/v1/puppies/{id}", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Detalhe de filhote")
	})
	mux.HandleFunc("PATCH /api/v1/puppies/{id}", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Edicao de filhote")
	})
	mux.HandleFunc("PATCH /api/v1/puppies/{id}/status", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Alteracao de status de filhote")
	})
	mux.HandleFunc("GET /api/v1/me/puppies", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Listagem dos meus filhotes")
	})
}
