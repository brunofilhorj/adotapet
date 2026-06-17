package chat

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
		webserver.HandleFunc("GET /ws/chat", Handler(), s.authenticate),
	}
}

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Chat WebSocket")
	}
}
