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

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func RegisterRoutes(
	mux *http.ServeMux,
	registerUsers inport.RegisterUserInputPort,
	loginUsers inport.LoginInputPort,
	refreshTokens inport.RefreshTokenInputPort,
) {
	mux.HandleFunc("POST /api/v1/auth/register", func(w http.ResponseWriter, r *http.Request) {
		handleRegister(w, r, registerUsers)
	})
	mux.HandleFunc("POST /api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, loginUsers)
	})
	mux.HandleFunc("POST /api/v1/auth/verify", func(w http.ResponseWriter, r *http.Request) {
		httperrors.NotImplemented(w, "Verificacao de conta")
	})
	mux.HandleFunc("POST /api/v1/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		handleRefresh(w, r, refreshTokens)
	})
}

func handleLogin(w http.ResponseWriter, r *http.Request, loginUsers inport.LoginInputPort) {
	var request loginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		httperrors.WriteJSON(w, http.StatusBadRequest, httperrors.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Payload de login invalido.",
		})
		return
	}

	tokens, err := loginUsers.Login(r.Context(), inport.LoginCommand{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		writeLoginError(w, err)
		return
	}

	httperrors.WriteJSON(w, http.StatusOK, map[string]any{
		"accessToken":      tokens.AccessToken,
		"refreshToken":     tokens.RefreshToken,
		"expiresIn":        tokens.ExpiresIn,
		"refreshExpiresIn": tokens.RefreshExpiresIn,
	})
}

func handleRefresh(w http.ResponseWriter, r *http.Request, refreshTokens inport.RefreshTokenInputPort) {
	var request refreshRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		httperrors.WriteJSON(w, http.StatusBadRequest, httperrors.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Payload de refresh invalido.",
		})
		return
	}

	tokens, err := refreshTokens.Refresh(r.Context(), inport.RefreshTokenCommand{
		RefreshToken: request.RefreshToken,
	})
	if err != nil {
		writeRefreshError(w, err)
		return
	}

	httperrors.WriteJSON(w, http.StatusOK, map[string]any{
		"accessToken":      tokens.AccessToken,
		"refreshToken":     tokens.RefreshToken,
		"expiresIn":        tokens.ExpiresIn,
		"refreshExpiresIn": tokens.RefreshExpiresIn,
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

func writeRefreshError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, authapp.ErrInvalidRefreshToken):
		httperrors.WriteJSON(w, http.StatusUnauthorized, httperrors.ErrorResponse{
			Code:    "INVALID_REFRESH_TOKEN",
			Message: "Refresh token invalido.",
		})
	case errors.Is(err, authapp.ErrRefreshTokenExpired):
		httperrors.WriteJSON(w, http.StatusUnauthorized, httperrors.ErrorResponse{
			Code:    "REFRESH_TOKEN_EXPIRED",
			Message: "Refresh token expirado ou revogado.",
		})
	case errors.Is(err, authapp.ErrAccountNotActive):
		httperrors.WriteJSON(w, http.StatusForbidden, httperrors.ErrorResponse{
			Code:    "ACCOUNT_NOT_VERIFIED",
			Message: "Conta ainda nao verificada.",
		})
	default:
		httperrors.WriteJSON(w, http.StatusInternalServerError, httperrors.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Erro interno ao renovar token.",
		})
	}
}

func writeLoginError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, authapp.ErrInvalidLoginCommand):
		httperrors.WriteJSON(w, http.StatusBadRequest, httperrors.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
	case errors.Is(err, authapp.ErrInvalidCredentials):
		httperrors.WriteJSON(w, http.StatusUnauthorized, httperrors.ErrorResponse{
			Code:    "INVALID_CREDENTIALS",
			Message: "Credenciais invalidas.",
		})
	case errors.Is(err, authapp.ErrAccountNotActive):
		httperrors.WriteJSON(w, http.StatusForbidden, httperrors.ErrorResponse{
			Code:    "ACCOUNT_NOT_VERIFIED",
			Message: "Conta ainda nao verificada.",
		})
	default:
		httperrors.WriteJSON(w, http.StatusInternalServerError, httperrors.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Erro interno ao autenticar usuario.",
		})
	}
}
