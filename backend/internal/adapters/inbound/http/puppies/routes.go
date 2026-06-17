package puppies

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
		webserver.HandleFunc("POST /api/v1/puppies", handleCreatePuppy, s.authenticate),
		webserver.HandleFunc("GET /api/v1/puppies/search", handleSearchPuppies),
		webserver.HandleFunc("GET /api/v1/puppies/{id}", handleGetPuppy),
		webserver.HandleFunc("PATCH /api/v1/puppies/{id}", handleUpdatePuppy, s.authenticate),
		webserver.HandleFunc("PATCH /api/v1/puppies/{id}/status", handleUpdatePuppyStatus, s.authenticate),
		webserver.HandleFunc("GET /api/v1/me/puppies", handleListMyPuppies, s.authenticate),
	}
}

func handleCreatePuppy(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Cadastro de filhote")
}

func handleSearchPuppies(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Busca geolocalizada de filhotes")
}

func handleGetPuppy(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Detalhe de filhote")
}

func handleUpdatePuppy(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Edicao de filhote")
}

func handleUpdatePuppyStatus(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Alteracao de status de filhote")
}

func handleListMyPuppies(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Listagem dos meus filhotes")
}
