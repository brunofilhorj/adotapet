package conversations

import (
	"net/http"

	httperrors "adotapet/internal/adapters/inbound/http/errors"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/conversations", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Inicio de conversa")
	})
	mux.HandleFunc("GET /api/v1/conversations", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Listagem de conversas")
	})
	mux.HandleFunc("GET /api/v1/conversations/{id}/messages", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Listagem de mensagens")
	})
	mux.HandleFunc("POST /api/v1/conversations/{id}/messages", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Envio de mensagem")
	})
}
