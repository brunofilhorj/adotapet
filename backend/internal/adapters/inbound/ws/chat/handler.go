package chat

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
		webserver.HandleFunc("GET /ws/chat", Handler()),
	}
}

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Chat WebSocket")
	}
}
