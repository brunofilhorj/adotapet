package media

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"adotapet/internal/adapters/inbound/http/webserver"
	authapp "adotapet/internal/app/auth"
)

func TestUploadURLRouteRequiresAccessToken(t *testing.T) {
	handler := webserver.New(nil, webserver.WithServices(NewService(fakeAccessTokenVerifier{})))

	request := httptest.NewRequest(http.MethodPost, "/api/v1/media/upload-url", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusUnauthorized, response.Body.String())
	}
}

type fakeAccessTokenVerifier struct{}

func (fakeAccessTokenVerifier) VerifyAccessToken(token string) (authapp.AccessTokenClaims, error) {
	return authapp.AccessTokenClaims{UserID: "user-1"}, nil
}
