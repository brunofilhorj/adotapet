package users

import (
	"encoding/json"
	"errors"
	"net/http"

	httperrors "adotapet/internal/adapters/inbound/http/errors"
	"adotapet/internal/adapters/inbound/http/middleware"
	"adotapet/internal/adapters/inbound/http/webserver"
	inport "adotapet/internal/app/port/in"
	usersapp "adotapet/internal/app/users"
	"adotapet/internal/domain/common"
)

type Service struct {
	authenticate webserver.Middleware
	profiles     profileService
}

type profileService interface {
	inport.GetMyProfileInputPort
	inport.UpdateProfileInputPort
}

type updateProfileRequest struct {
	Name      *string        `json:"name"`
	Phone     *string        `json:"phone"`
	City      *string        `json:"city"`
	State     *string        `json:"state"`
	Location  *locationInput `json:"location"`
	AvatarURL *string        `json:"avatarUrl"`
	Bio       *string        `json:"bio"`
}

type locationInput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type profileResponse struct {
	UserID    string           `json:"userId"`
	Email     string           `json:"email"`
	Role      string           `json:"role"`
	Status    string           `json:"status"`
	Name      string           `json:"name"`
	Phone     *string          `json:"phone,omitempty"`
	City      string           `json:"city"`
	State     string           `json:"state"`
	Location  *common.GeoPoint `json:"location,omitempty"`
	AvatarURL *string          `json:"avatarUrl,omitempty"`
	Bio       *string          `json:"bio,omitempty"`
}

func NewService(accessTokenVerifier middleware.AccessTokenVerifier, profiles profileService) Service {
	return Service{
		authenticate: middleware.Authenticate(accessTokenVerifier),
		profiles:     profiles,
	}
}

func (s Service) Routes() []webserver.Route {
	return []webserver.Route{
		webserver.HandleFunc("GET /api/v1/me", func(w http.ResponseWriter, r *http.Request) {
			handleGetMe(w, r, s.profiles)
		}, s.authenticate),
		webserver.HandleFunc("PATCH /api/v1/me/profile", func(w http.ResponseWriter, r *http.Request) {
			handleUpdateProfile(w, r, s.profiles)
		}, s.authenticate),
	}
}

func handleGetMe(w http.ResponseWriter, r *http.Request, profiles inport.GetMyProfileInputPort) {
	claims, ok := middleware.AuthenticatedUser(r.Context())
	if !ok {
		httperrors.WriteJSON(w, http.StatusUnauthorized, httperrors.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Usuario autenticado ausente.",
		})
		return
	}

	profile, err := profiles.Get(r.Context(), inport.GetMyProfileQuery{UserID: claims.UserID})
	if err != nil {
		writeProfileError(w, err)
		return
	}

	httperrors.WriteJSON(w, http.StatusOK, toProfileResponse(profile))
}

func handleUpdateProfile(w http.ResponseWriter, r *http.Request, profiles inport.UpdateProfileInputPort) {
	claims, ok := middleware.AuthenticatedUser(r.Context())
	if !ok {
		httperrors.WriteJSON(w, http.StatusUnauthorized, httperrors.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Usuario autenticado ausente.",
		})
		return
	}

	var request updateProfileRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		httperrors.WriteJSON(w, http.StatusBadRequest, httperrors.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Payload de perfil invalido.",
		})
		return
	}

	profile, err := profiles.Update(r.Context(), inport.UpdateProfileCommand{
		UserID:    claims.UserID,
		Name:      request.Name,
		Phone:     request.Phone,
		City:      request.City,
		State:     request.State,
		Location:  toGeoPoint(request.Location),
		AvatarURL: request.AvatarURL,
		Bio:       request.Bio,
	})
	if err != nil {
		writeProfileError(w, err)
		return
	}

	httperrors.WriteJSON(w, http.StatusOK, toProfileResponse(profile))
}

func toGeoPoint(input *locationInput) *common.GeoPoint {
	if input == nil {
		return nil
	}
	return &common.GeoPoint{
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
	}
}

func toProfileResponse(profile inport.MyProfile) profileResponse {
	return profileResponse{
		UserID:    profile.UserID,
		Email:     profile.Email,
		Role:      profile.Role,
		Status:    profile.Status,
		Name:      profile.Name,
		Phone:     profile.Phone,
		City:      profile.City,
		State:     profile.State,
		Location:  profile.Location,
		AvatarURL: profile.AvatarURL,
		Bio:       profile.Bio,
	}
}

func writeProfileError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, usersapp.ErrUserNotFound), errors.Is(err, usersapp.ErrProfileNotFound):
		httperrors.WriteJSON(w, http.StatusNotFound, httperrors.ErrorResponse{
			Code:    "PROFILE_NOT_FOUND",
			Message: "Perfil nao encontrado.",
		})
	case errors.Is(err, usersapp.ErrInvalidProfileCommand):
		httperrors.WriteJSON(w, http.StatusBadRequest, httperrors.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
	default:
		httperrors.WriteJSON(w, http.StatusInternalServerError, httperrors.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Erro interno ao processar perfil.",
		})
	}
}
