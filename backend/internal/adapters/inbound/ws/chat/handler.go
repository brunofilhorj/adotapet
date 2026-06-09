package chat

import (
	"net/http"

	httperrors "adotapet/internal/adapters/inbound/http/errors"
)

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Chat WebSocket")
	}
}
