package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	httperrors "adotapet/internal/adapters/inbound/http/errors"
	authapp "adotapet/internal/app/auth"
	inport "adotapet/internal/app/port/in"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Name     string `json:"name"`
	City     string `json:"city"`
	State    string `json:"state"`
}

func RegisterRoutes(mux *http.ServeMux, registerUsers inport.RegisterUserInputPort) {
	mux.HandleFunc("POST /api/v1/auth/register", func(w http.ResponseWriter, r *http.Request) {
		handleRegister(w, r, registerUsers)
	})
	mux.HandleFunc("POST /api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Login")
	})
	mux.HandleFunc("POST /api/v1/auth/verify", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Verificacao de conta")
	})
	mux.HandleFunc("POST /api/v1/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Refresh token")
	})
}

func handleRegister(w http.ResponseWriter, r *http.Request, registerUsers inport.RegisterUserInputPort) {
	var request registerRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		httperrors.WriteJSON(w, http.StatusBadRequest, httperrors.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Payload de cadastro invalido.",
		})
		return
	}

	registered, err := registerUsers.Register(r.Context(), inport.RegisterUserCommand{
		Email:    request.Email,
		Password: request.Password,
		Role:     request.Role,
		Name:     request.Name,
		City:     request.City,
		State:    request.State,
	})
	if err != nil {
		writeRegisterError(w, err)
		return
	}

	httperrors.WriteJSON(w, http.StatusCreated, map[string]string{
		"userId": registered.UserID,
		"status": registered.Status,
	})
}

func writeRegisterError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, authapp.ErrEmailAlreadyRegistered):
		httperrors.WriteJSON(w, http.StatusConflict, httperrors.ErrorResponse{
			Code:    "EMAIL_ALREADY_REGISTERED",
			Message: "E-mail ja cadastrado.",
		})
	case errors.Is(err, authapp.ErrInvalidRegisterCommand):
		httperrors.WriteJSON(w, http.StatusBadRequest, httperrors.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
	default:
		httperrors.WriteJSON(w, http.StatusInternalServerError, httperrors.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Erro interno ao cadastrar usuario.",
		})
	}
}
