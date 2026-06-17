package conversations

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
		webserver.HandleFunc("POST /api/v1/conversations", handleCreateConversation, s.authenticate),
		webserver.HandleFunc("GET /api/v1/conversations", handleListConversations, s.authenticate),
		webserver.HandleFunc("GET /api/v1/conversations/{id}/messages", handleListMessages, s.authenticate),
		webserver.HandleFunc("POST /api/v1/conversations/{id}/messages", handleCreateMessage, s.authenticate),
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
