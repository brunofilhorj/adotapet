package favorites

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
		webserver.HandleFunc("POST /api/v1/puppies/{id}/favorite", handleAddFavorite),
		webserver.HandleFunc("DELETE /api/v1/puppies/{id}/favorite", handleRemoveFavorite),
		webserver.HandleFunc("GET /api/v1/me/favorites", handleListFavorites),
		webserver.HandleFunc("POST /api/v1/me/saved-searches", handleCreateSavedSearch),
		webserver.HandleFunc("GET /api/v1/me/saved-searches", handleListSavedSearches),
		webserver.HandleFunc("DELETE /api/v1/me/saved-searches/{id}", handleDeleteSavedSearch),
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
