package puppies

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
	"adotapet/internal/domain/common"
)

func TestProtectedPuppyRoutesRequireAccessToken(t *testing.T) {
	handler := webserver.New(nil, webserver.WithServices(NewService(fakeAccessTokenVerifier{}, &fakeListingService{})))

	assertStatus(t, handler, http.MethodPost, "/api/v1/puppies", http.StatusUnauthorized)
	assertStatus(t, handler, http.MethodPatch, "/api/v1/puppies/puppy-1", http.StatusUnauthorized)
	assertStatus(t, handler, http.MethodPatch, "/api/v1/puppies/puppy-1/status", http.StatusUnauthorized)
	assertStatus(t, handler, http.MethodGet, "/api/v1/me/puppies", http.StatusUnauthorized)
}

func TestPublicPuppyRoutesStayPublic(t *testing.T) {
	handler := webserver.New(nil, webserver.WithServices(NewService(fakeAccessTokenVerifier{}, &fakeListingService{
		detail: puppyDetails("puppy-1"),
	})))

	assertStatus(t, handler, http.MethodGet, "/api/v1/puppies/search", http.StatusNotImplemented)
	assertStatus(t, handler, http.MethodGet, "/api/v1/puppies/puppy-1", http.StatusOK)
}

func TestCreatePuppyUsesAuthenticatedUser(t *testing.T) {
	listings := &fakeListingService{detail: puppyDetails("puppy-1")}
	handler := webserver.New(nil, webserver.WithServices(NewService(fakeAccessTokenVerifier{}, listings)))

	body := bytes.NewBufferString(`{
		"name":"Luna",
		"breed":"SRD",
		"species":"DOG",
		"ageMonths":3,
		"size":"SMALL",
		"sex":"FEMALE",
		"description":"Muito docil",
		"location":{"latitude":-23.55,"longitude":-46.63},
		"city":"Sao Paulo",
		"state":"SP"
	}`)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/puppies", body)
	request.Header.Set("Authorization", "Bearer valid-token")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusCreated, response.Body.String())
	}
	if listings.created.OwnerID != "user-1" || listings.created.OwnerRole != "DONOR" {
		t.Fatalf("created = %+v, want authenticated owner", listings.created)
	}
	if !strings.Contains(response.Body.String(), `"id":"puppy-1"`) {
		t.Fatalf("body = %s, want created puppy", response.Body.String())
	}
}

func TestListMyPuppiesUsesAuthenticatedUser(t *testing.T) {
	listings := &fakeListingService{
		page: common.Page[inport.PuppyDetails]{
			Items:         []inport.PuppyDetails{puppyDetails("puppy-1")},
			Page:          2,
			Size:          10,
			TotalElements: 1,
			TotalPages:    1,
		},
	}
	handler := webserver.New(nil, webserver.WithServices(NewService(fakeAccessTokenVerifier{}, listings)))

	request := httptest.NewRequest(http.MethodGet, "/api/v1/me/puppies?page=2&size=10", nil)
	request.Header.Set("Authorization", "Bearer valid-token")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	if listings.listOwnerID != "user-1" || listings.listPage.Page != 2 || listings.listPage.Size != 10 {
		t.Fatalf("list owner/page = %q %+v, want authenticated owner and query page", listings.listOwnerID, listings.listPage)
	}
}

type fakeAccessTokenVerifier struct{}

func (fakeAccessTokenVerifier) VerifyAccessToken(token string) (authapp.AccessTokenClaims, error) {
	return authapp.AccessTokenClaims{UserID: "user-1", Role: "DONOR", Status: "ACTIVE"}, nil
}

type fakeListingService struct {
	detail      inport.PuppyDetails
	page        common.Page[inport.PuppyDetails]
	created     inport.CreatePuppyCommand
	listOwnerID string
	listPage    common.PageRequest
}

func (s *fakeListingService) Create(ctx context.Context, cmd inport.CreatePuppyCommand) (inport.PuppyDetails, error) {
	s.created = cmd
	return s.detail, nil
}

func (s *fakeListingService) Get(ctx context.Context, query inport.GetPuppyQuery) (inport.PuppyDetails, error) {
	return s.detail, nil
}

func (s *fakeListingService) ListMine(ctx context.Context, query inport.ListMyPuppiesQuery) (common.Page[inport.PuppyDetails], error) {
	s.listOwnerID = query.OwnerID
	s.listPage = query.Page
	return s.page, nil
}

func puppyDetails(id string) inport.PuppyDetails {
	return inport.PuppyDetails{
		ID:          id,
		OwnerID:     "user-1",
		Name:        "Luna",
		Species:     "DOG",
		AgeMonths:   3,
		Size:        "SMALL",
		Sex:         "FEMALE",
		Description: "Muito docil",
		Location:    common.GeoPoint{Latitude: -23.55, Longitude: -46.63},
		City:        "Sao Paulo",
		State:       "SP",
		Status:      "AVAILABLE",
		CreatedAt:   "2026-06-17T00:00:00Z",
		UpdatedAt:   "2026-06-17T00:00:00Z",
	}
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
