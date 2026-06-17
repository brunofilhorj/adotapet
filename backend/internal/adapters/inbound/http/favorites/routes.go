package favorites

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
		webserver.HandleFunc("POST /api/v1/puppies/{id}/favorite", handleAddFavorite, s.authenticate),
		webserver.HandleFunc("DELETE /api/v1/puppies/{id}/favorite", handleRemoveFavorite, s.authenticate),
		webserver.HandleFunc("GET /api/v1/me/favorites", handleListFavorites, s.authenticate),
		webserver.HandleFunc("POST /api/v1/me/saved-searches", handleCreateSavedSearch, s.authenticate),
		webserver.HandleFunc("GET /api/v1/me/saved-searches", handleListSavedSearches, s.authenticate),
		webserver.HandleFunc("DELETE /api/v1/me/saved-searches/{id}", handleDeleteSavedSearch, s.authenticate),
	}
}

func handleAddFavorite(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Adicionar favorito")
}

func handleRemoveFavorite(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Remover favorito")
}

func handleListFavorites(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Listagem de favoritos")
}

func handleCreateSavedSearch(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Criacao de busca salva")
}

func handleListSavedSearches(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Listagem de buscas salvas")
}

func handleDeleteSavedSearch(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Remocao de busca salva")
}
