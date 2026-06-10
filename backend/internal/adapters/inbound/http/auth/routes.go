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
	Email               string `json:"email"`
	Password            string `json:"password"`
	Role                string `json:"role"`
	Name                string `json:"name"`
	City                string `json:"city"`
	State               string `json:"state"`
	Phone               string `json:"phone"`
	VerificationChannel string `json:"verificationChannel"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type verifyRequest struct {
	Email       string `json:"email"`
	Channel     string `json:"channel"`
	Destination string `json:"destination"`
	Code        string `json:"code"`
}

type resendVerificationRequest struct {
	Email       string `json:"email"`
	Channel     string `json:"channel"`
	Destination string `json:"destination"`
}

func RegisterRoutes(
	mux *http.ServeMux,
	registerUsers inport.RegisterUserInputPort,
	loginUsers inport.LoginInputPort,
	refreshTokens inport.RefreshTokenInputPort,
	verifyAccounts inport.VerifyAccountInputPort,
	resendVerification inport.ResendVerificationInputPort,
) {
	mux.HandleFunc("POST /api/v1/auth/register", func(w http.ResponseWriter, r *http.Request) {
		handleRegister(w, r, registerUsers)
	})
	mux.HandleFunc("POST /api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, loginUsers)
	})
	mux.HandleFunc("POST /api/v1/auth/verify", func(w http.ResponseWriter, r *http.Request) {
		handleVerify(w, r, verifyAccounts)
	})
	mux.HandleFunc("POST /api/v1/auth/resend-verification", func(w http.ResponseWriter, r *http.Request) {
		handleResendVerification(w, r, resendVerification)
	})
	mux.HandleFunc("POST /api/v1/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		handleRefresh(w, r, refreshTokens)
	})
}

func handleVerify(w http.ResponseWriter, r *http.Request, verifyAccounts inport.VerifyAccountInputPort) {
	var request verifyRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		httperrors.WriteJSON(w, http.StatusBadRequest, httperrors.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Payload de verificacao invalido.",
		})
		return
	}

	verified, err := verifyAccounts.Verify(r.Context(), inport.VerifyAccountCommand{
		Email:       request.Email,
		Channel:     request.Channel,
		Destination: request.Destination,
		Code:        request.Code,
	})
	if err != nil {
		writeVerificationError(w, err)
		return
	}

	httperrors.WriteJSON(w, http.StatusOK, map[string]string{
		"userId": verified.UserID,
		"status": verified.Status,
	})
}

func handleResendVerification(w http.ResponseWriter, r *http.Request, resendVerification inport.ResendVerificationInputPort) {
	var request resendVerificationRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		httperrors.WriteJSON(w, http.StatusBadRequest, httperrors.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Payload de reenvio invalido.",
		})
		return
	}

	result, err := resendVerification.Resend(r.Context(), inport.ResendVerificationCommand{
		Email:       request.Email,
		Channel:     request.Channel,
		Destination: request.Destination,
	})
	if err != nil {
		writeVerificationError(w, err)
		return
	}

	httperrors.WriteJSON(w, http.StatusAccepted, map[string]string{
		"userId":  result.UserID,
		"channel": result.Channel,
		"target":  result.Target,
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
		Email:               request.Email,
		Password:            request.Password,
		Role:                request.Role,
		Name:                request.Name,
		City:                request.City,
		State:               request.State,
		Phone:               request.Phone,
		VerificationChannel: request.VerificationChannel,
	})
	if err != nil {
		writeRegisterError(w, err)
		return
	}

	httperrors.WriteJSON(w, http.StatusCreated, map[string]string{
		"userId":              registered.UserID,
		"status":              registered.Status,
		"verificationChannel": registered.VerificationChannel,
		"verificationTarget":  registered.VerificationTarget,
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

func writeVerificationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, authapp.ErrInvalidVerificationCommand):
		httperrors.WriteJSON(w, http.StatusBadRequest, httperrors.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
	case errors.Is(err, authapp.ErrInvalidVerificationChannel):
		httperrors.WriteJSON(w, http.StatusBadRequest, httperrors.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Canal de verificacao invalido.",
		})
	case errors.Is(err, authapp.ErrVerificationDestinationMissing):
		httperrors.WriteJSON(w, http.StatusBadRequest, httperrors.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Destino de verificacao obrigatorio para o canal.",
		})
	case errors.Is(err, authapp.ErrVerificationCodeInvalid):
		httperrors.WriteJSON(w, http.StatusUnauthorized, httperrors.ErrorResponse{
			Code:    "INVALID_VERIFICATION_CODE",
			Message: "Codigo de verificacao invalido.",
		})
	case errors.Is(err, authapp.ErrVerificationCodeExpired):
		httperrors.WriteJSON(w, http.StatusUnauthorized, httperrors.ErrorResponse{
			Code:    "VERIFICATION_CODE_EXPIRED",
			Message: "Codigo de verificacao expirado.",
		})
	case errors.Is(err, authapp.ErrAccountAlreadyActive):
		httperrors.WriteJSON(w, http.StatusConflict, httperrors.ErrorResponse{
			Code:    "ACCOUNT_ALREADY_ACTIVE",
			Message: "Conta ja esta ativa.",
		})
	default:
		httperrors.WriteJSON(w, http.StatusInternalServerError, httperrors.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Erro interno ao verificar conta.",
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
