package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	httperrors "adotapet/internal/adapters/inbound/http/errors"
	authapp "adotapet/internal/app/auth"
)

type AccessTokenVerifier interface {
	VerifyAccessToken(token string) (authapp.AccessTokenClaims, error)
}

type authenticatedUserKey struct{}

func Authenticate(verifier AccessTokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := bearerToken(r.Header.Get("Authorization"))
			if !ok {
				httperrors.WriteJSON(w, http.StatusUnauthorized, httperrors.ErrorResponse{
					Code:    "UNAUTHORIZED",
					Message: "Token de acesso ausente.",
				})
				return
			}

			claims, err := verifier.VerifyAccessToken(token)
			if err != nil {
				code := "INVALID_ACCESS_TOKEN"
				message := "Token de acesso invalido."
				if errors.Is(err, authapp.ErrAccessTokenExpired) {
					code = "ACCESS_TOKEN_EXPIRED"
					message = "Token de acesso expirado."
				}

				httperrors.WriteJSON(w, http.StatusUnauthorized, httperrors.ErrorResponse{
					Code:    code,
					Message: message,
				})
				return
			}

			ctx := context.WithValue(r.Context(), authenticatedUserKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthenticatedUser(ctx context.Context) (authapp.AccessTokenClaims, bool) {
	claims, ok := ctx.Value(authenticatedUserKey{}).(authapp.AccessTokenClaims)
	return claims, ok
}

func bearerToken(header string) (string, bool) {
	prefix := "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	return token, token != ""
}
