package media

import (
	"net/http"

	httperrors "adotapet/internal/adapters/inbound/http/errors"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/media/upload-url", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "URL pre-assinada de upload")
	})
}
