package users

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"adotapet/internal/adapters/inbound/http/webserver"
	authapp "adotapet/internal/app/auth"
	inport "adotapet/internal/app/port/in"
)

func TestGetMeReturnsAuthenticatedProfile(t *testing.T) {
	profiles := &fakeProfileService{
		profile: inport.MyProfile{
			UserID: "user-1",
			Email:  "maria@example.com",
			Role:   "ADOPTER",
			Status: "ACTIVE",
			Name:   "Maria Souza",
			City:   "Sao Paulo",
			State:  "SP",
		},
	}
	handler := webserver.New(nil, webserver.WithServices(NewService(fakeAccessTokenVerifier{}, profiles)))

	request := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	request.Header.Set("Authorization", "Bearer valid-token")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	if profiles.getUserID != "user-1" {
		t.Fatalf("getUserID = %q, want user-1", profiles.getUserID)
	}
	if !strings.Contains(response.Body.String(), `"email":"maria@example.com"`) {
		t.Fatalf("body = %s, want profile json", response.Body.String())
	}
}

func TestUpdateProfileUsesAuthenticatedUser(t *testing.T) {
	profiles := &fakeProfileService{
		profile: inport.MyProfile{
			UserID: "user-1",
			Email:  "maria@example.com",
			Role:   "ADOPTER",
			Status: "ACTIVE",
			Name:   "Maria Oliveira",
			City:   "Rio de Janeiro",
			State:  "RJ",
		},
	}
	handler := webserver.New(nil, webserver.WithServices(NewService(fakeAccessTokenVerifier{}, profiles)))

	body := bytes.NewBufferString(`{"name":"Maria Oliveira","city":"Rio de Janeiro","state":"RJ"}`)
	request := httptest.NewRequest(http.MethodPatch, "/api/v1/me/profile", body)
	request.Header.Set("Authorization", "Bearer valid-token")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	if profiles.updateUserID != "user-1" {
		t.Fatalf("updateUserID = %q, want user-1", profiles.updateUserID)
	}
	if profiles.updateName == nil || *profiles.updateName != "Maria Oliveira" {
		t.Fatalf("updateName = %v, want Maria Oliveira", profiles.updateName)
	}
}

type fakeAccessTokenVerifier struct{}

func (fakeAccessTokenVerifier) VerifyAccessToken(token string) (authapp.AccessTokenClaims, error) {
	return authapp.AccessTokenClaims{
		UserID: "user-1",
		Email:  "maria@example.com",
		Role:   "ADOPTER",
		Status: "ACTIVE",
	}, nil
}

type fakeProfileService struct {
	profile      inport.MyProfile
	getUserID    string
	updateUserID string
	updateName   *string
}

func (s *fakeProfileService) Get(ctx context.Context, query inport.GetMyProfileQuery) (inport.MyProfile, error) {
	s.getUserID = query.UserID
	return s.profile, nil
}

func (s *fakeProfileService) Update(ctx context.Context, cmd inport.UpdateProfileCommand) (inport.MyProfile, error) {
	s.updateUserID = cmd.UserID
	s.updateName = cmd.Name
	return s.profile, nil
}
