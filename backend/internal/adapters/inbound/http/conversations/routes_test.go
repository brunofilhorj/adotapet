package conversations

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"adotapet/internal/adapters/inbound/http/webserver"
	authapp "adotapet/internal/app/auth"
)

func TestConversationRoutesRequireAccessToken(t *testing.T) {
	handler := webserver.New(nil, webserver.WithServices(NewService(fakeAccessTokenVerifier{})))

	assertStatus(t, handler, http.MethodPost, "/api/v1/conversations", http.StatusUnauthorized)
	assertStatus(t, handler, http.MethodGet, "/api/v1/conversations", http.StatusUnauthorized)
	assertStatus(t, handler, http.MethodGet, "/api/v1/conversations/conversation-1/messages", http.StatusUnauthorized)
	assertStatus(t, handler, http.MethodPost, "/api/v1/conversations/conversation-1/messages", http.StatusUnauthorized)
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
