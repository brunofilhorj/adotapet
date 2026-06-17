package conversations

import (
	"net/http"

	httperrors "adotapet/internal/adapters/inbound/http/errors"
	"adotapet/internal/adapters/inbound/http/webserver"
)

type Service struct{}

func NewService() Service {
	return Service{}
}

func (Service) Routes() []webserver.Route {
	return []webserver.Route{
		webserver.HandleFunc("POST /api/v1/conversations", handleCreateConversation),
		webserver.HandleFunc("GET /api/v1/conversations", handleListConversations),
		webserver.HandleFunc("GET /api/v1/conversations/{id}/messages", handleListMessages),
		webserver.HandleFunc("POST /api/v1/conversations/{id}/messages", handleCreateMessage),
	}
}

func handleCreateConversation(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Inicio de conversa")
}

func handleListConversations(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Listagem de conversas")
}

func handleListMessages(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Listagem de mensagens")
}

func handleCreateMessage(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Envio de mensagem")
}
