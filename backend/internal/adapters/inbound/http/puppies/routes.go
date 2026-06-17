package puppies

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	httperrors "adotapet/internal/adapters/inbound/http/errors"
	"adotapet/internal/adapters/inbound/http/middleware"
	"adotapet/internal/adapters/inbound/http/webserver"
	inport "adotapet/internal/app/port/in"
	puppiesapp "adotapet/internal/app/puppies"
	"adotapet/internal/domain/common"
)

type Service struct {
	authenticate webserver.Middleware
	listings     listingService
}

type listingService interface {
	inport.CreatePuppyInputPort
	inport.GetPuppyInputPort
	inport.ListMyPuppiesInputPort
}

type createPuppyRequest struct {
	Name        string          `json:"name"`
	Breed       *string         `json:"breed"`
	Species     string          `json:"species"`
	AgeMonths   int             `json:"ageMonths"`
	Size        string          `json:"size"`
	Sex         string          `json:"sex"`
	Description string          `json:"description"`
	Location    common.GeoPoint `json:"location"`
	City        string          `json:"city"`
	State       string          `json:"state"`
}

type puppyResponse struct {
	ID          string          `json:"id"`
	OwnerID     string          `json:"ownerId"`
	Name        string          `json:"name"`
	Breed       *string         `json:"breed,omitempty"`
	Species     string          `json:"species"`
	AgeMonths   int             `json:"ageMonths"`
	Size        string          `json:"size"`
	Sex         string          `json:"sex"`
	Description string          `json:"description"`
	Location    common.GeoPoint `json:"location"`
	City        string          `json:"city"`
	State       string          `json:"state"`
	Status      string          `json:"status"`
	AdoptedAt   *string         `json:"adoptedAt,omitempty"`
	CreatedAt   string          `json:"createdAt"`
	UpdatedAt   string          `json:"updatedAt"`
}

type pageResponse[T any] struct {
	Items         []T   `json:"items"`
	Page          int   `json:"page"`
	Size          int   `json:"size"`
	TotalElements int64 `json:"totalElements"`
	TotalPages    int   `json:"totalPages"`
}

func NewService(accessTokenVerifier middleware.AccessTokenVerifier, listings listingService) Service {
	return Service{
		authenticate: middleware.Authenticate(accessTokenVerifier),
		listings:     listings,
	}
}

func (s Service) Routes() []webserver.Route {
	return []webserver.Route{
		webserver.HandleFunc("POST /api/v1/puppies", func(w http.ResponseWriter, r *http.Request) {
			handleCreatePuppy(w, r, s.listings)
		}, s.authenticate),
		webserver.HandleFunc("GET /api/v1/puppies/search", handleSearchPuppies),
		webserver.HandleFunc("GET /api/v1/puppies/{id}", func(w http.ResponseWriter, r *http.Request) {
			handleGetPuppy(w, r, s.listings)
		}),
		webserver.HandleFunc("PATCH /api/v1/puppies/{id}", handleUpdatePuppy, s.authenticate),
		webserver.HandleFunc("PATCH /api/v1/puppies/{id}/status", handleUpdatePuppyStatus, s.authenticate),
		webserver.HandleFunc("GET /api/v1/me/puppies", func(w http.ResponseWriter, r *http.Request) {
			handleListMyPuppies(w, r, s.listings)
		}, s.authenticate),
	}
}

func handleCreatePuppy(w http.ResponseWriter, r *http.Request, listings inport.CreatePuppyInputPort) {
	claims, ok := middleware.AuthenticatedUser(r.Context())
	if !ok {
		httperrors.WriteJSON(w, http.StatusUnauthorized, httperrors.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Usuario autenticado ausente.",
		})
		return
	}

	var request createPuppyRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		httperrors.WriteJSON(w, http.StatusBadRequest, httperrors.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Payload de filhote invalido.",
		})
		return
	}

	created, err := listings.Create(r.Context(), inport.CreatePuppyCommand{
		OwnerID:     claims.UserID,
		OwnerRole:   claims.Role,
		Name:        request.Name,
		Breed:       request.Breed,
		Species:     request.Species,
		AgeMonths:   request.AgeMonths,
		Size:        request.Size,
		Sex:         request.Sex,
		Description: request.Description,
		Location:    request.Location,
		City:        request.City,
		State:       request.State,
	})
	if err != nil {
		writePuppyError(w, err)
		return
	}

	httperrors.WriteJSON(w, http.StatusCreated, toPuppyResponse(created))
}

func handleSearchPuppies(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Busca geolocalizada de filhotes")
}

func handleGetPuppy(w http.ResponseWriter, r *http.Request, listings inport.GetPuppyInputPort) {
	found, err := listings.Get(r.Context(), inport.GetPuppyQuery{PuppyID: r.PathValue("id")})
	if err != nil {
		writePuppyError(w, err)
		return
	}

	httperrors.WriteJSON(w, http.StatusOK, toPuppyResponse(found))
}

func handleUpdatePuppy(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Edicao de filhote")
}

func handleUpdatePuppyStatus(w http.ResponseWriter, r *http.Request) {
	httperrors.NotImplemented(w, "Alteracao de status de filhote")
}

func handleListMyPuppies(w http.ResponseWriter, r *http.Request, listings inport.ListMyPuppiesInputPort) {
	claims, ok := middleware.AuthenticatedUser(r.Context())
	if !ok {
		httperrors.WriteJSON(w, http.StatusUnauthorized, httperrors.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Usuario autenticado ausente.",
		})
		return
	}

	page, err := listings.ListMine(r.Context(), inport.ListMyPuppiesQuery{
		OwnerID: claims.UserID,
		Page: common.PageRequest{
			Page: intQuery(r, "page", 1),
			Size: intQuery(r, "size", 20),
		},
	})
	if err != nil {
		writePuppyError(w, err)
		return
	}

	items := make([]puppyResponse, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, toPuppyResponse(item))
	}

	httperrors.WriteJSON(w, http.StatusOK, pageResponse[puppyResponse]{
		Items:         items,
		Page:          page.Page,
		Size:          page.Size,
		TotalElements: page.TotalElements,
		TotalPages:    page.TotalPages,
	})
}

func intQuery(r *http.Request, name string, fallback int) int {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}

func toPuppyResponse(details inport.PuppyDetails) puppyResponse {
	return puppyResponse{
		ID:          details.ID,
		OwnerID:     details.OwnerID,
		Name:        details.Name,
		Breed:       details.Breed,
		Species:     details.Species,
		AgeMonths:   details.AgeMonths,
		Size:        details.Size,
		Sex:         details.Sex,
		Description: details.Description,
		Location:    details.Location,
		City:        details.City,
		State:       details.State,
		Status:      details.Status,
		AdoptedAt:   details.AdoptedAt,
		CreatedAt:   details.CreatedAt,
		UpdatedAt:   details.UpdatedAt,
	}
}

func writePuppyError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, puppiesapp.ErrPuppyNotFound):
		httperrors.WriteJSON(w, http.StatusNotFound, httperrors.ErrorResponse{
			Code:    "PUPPY_NOT_FOUND",
			Message: "Filhote nao encontrado.",
		})
	case errors.Is(err, puppiesapp.ErrInvalidPuppyCommand):
		httperrors.WriteJSON(w, http.StatusBadRequest, httperrors.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
	case errors.Is(err, puppiesapp.ErrPuppyForbidden):
		httperrors.WriteJSON(w, http.StatusForbidden, httperrors.ErrorResponse{
			Code:    "FORBIDDEN",
			Message: "Usuario nao autorizado para criar anuncios.",
		})
	default:
		httperrors.WriteJSON(w, http.StatusInternalServerError, httperrors.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Erro interno ao processar filhote.",
		})
	}
}
