package media

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
		webserver.HandleFunc("POST /api/v1/media/upload-url", handleUploadURL),
	}
}

func handleUploadURL(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "URL pre-assinada de upload")
}
