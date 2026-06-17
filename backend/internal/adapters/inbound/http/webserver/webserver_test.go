package webserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRegistersRoutesFromRoutesAndServices(t *testing.T) {
	handler := New(
		nil,
		WithRoutes(HandleFunc("GET /health", writeText("health"))),
		WithServices(testService{}),
	)

	assertResponse(t, handler, "GET", "/health", http.StatusOK, "health")
	assertResponse(t, handler, "GET", "/service", http.StatusOK, "service")
}

func TestNewAppliesGlobalAndRouteMiddlewares(t *testing.T) {
	handler := New(
		nil,
		WithRoutes(HandleFunc(
			"GET /middleware",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(r.Header.Get("X-Trace")))
			},
			appendHeader("route"),
		)),
		WithMiddlewares(appendHeader("global")),
	)

	assertResponse(t, handler, "GET", "/middleware", http.StatusOK, "global,route")
}

type testService struct{}

func (testService) Routes() []Route {
	return []Route{
		HandleFunc("GET /service", writeText("service")),
	}
}

func writeText(value string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(value))
	}
}

func appendHeader(value string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			current := r.Header.Get("X-Trace")
			if current == "" {
				r.Header.Set("X-Trace", value)
			} else {
				r.Header.Set("X-Trace", current+","+value)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func assertResponse(t *testing.T, handler http.Handler, method string, path string, status int, body string) {
	t.Helper()

	request := httptest.NewRequest(method, path, nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != status {
		t.Fatalf("status = %d, want %d", response.Code, status)
	}
	if response.Body.String() != body {
		t.Fatalf("body = %q, want %q", response.Body.String(), body)
	}
}
