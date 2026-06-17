package puppies

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"adotapet/internal/adapters/inbound/http/webserver"
	authapp "adotapet/internal/app/auth"
)

func TestProtectedPuppyRoutesRequireAccessToken(t *testing.T) {
	handler := webserver.New(nil, webserver.WithServices(NewService(fakeAccessTokenVerifier{})))

	assertStatus(t, handler, http.MethodPost, "/api/v1/puppies", http.StatusUnauthorized)
	assertStatus(t, handler, http.MethodPatch, "/api/v1/puppies/puppy-1", http.StatusUnauthorized)
	assertStatus(t, handler, http.MethodPatch, "/api/v1/puppies/puppy-1/status", http.StatusUnauthorized)
	assertStatus(t, handler, http.MethodGet, "/api/v1/me/puppies", http.StatusUnauthorized)
}

func TestPublicPuppyRoutesStayPublic(t *testing.T) {
	handler := webserver.New(nil, webserver.WithServices(NewService(fakeAccessTokenVerifier{})))

	assertStatus(t, handler, http.MethodGet, "/api/v1/puppies/search", http.StatusNotImplemented)
	assertStatus(t, handler, http.MethodGet, "/api/v1/puppies/puppy-1", http.StatusNotImplemented)
}

type fakeAccessTokenVerifier struct{}

func (fakeAccessTokenVerifier) VerifyAccessToken(token string) (authapp.AccessTokenClaims, error) {
	return authapp.AccessTokenClaims{UserID: "user-1"}, nil
}

func assertStatus(t *testing.T, handler http.Handler, method string, path string, status int) {
	t.Helper()

	request := httptest.NewRequest(method, path, nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != status {
		t.Fatalf("%s %s status = %d, want %d; body = %s", method, path, response.Code, status, response.Body.String())
	}
}
