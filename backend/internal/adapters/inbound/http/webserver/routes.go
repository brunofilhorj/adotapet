package webserver

import "net/http"

// Route describes a single HTTP route registered on a ServeMux.
type Route struct {
	pattern     string
	handler     http.Handler
	middlewares []Middleware
}

// RouteProvider provides the routes owned by a resource.
type RouteProvider interface {
	Routes() []Route
}

// Handle creates a route with optional middlewares applied only to this route.
func Handle(pattern string, handler http.Handler, middlewares ...Middleware) Route {
	return Route{
		pattern:     pattern,
		handler:     handler,
		middlewares: middlewares,
	}
}

// HandleFunc creates a route from a handler function with optional route middlewares.
func HandleFunc(pattern string, handlerFunc http.HandlerFunc, middlewares ...Middleware) Route {
	return Handle(pattern, handlerFunc, middlewares...)
}

// WithRoutes registers routes on a new ServeMux and uses it as the root handler.
func WithRoutes(routes ...Route) Option {
	return func(config *config) {
		config.routes = append(config.routes, routes...)
	}
}

// WithServices registers routes exposed by resource route providers.
func WithServices(providers ...RouteProvider) Option {
	return func(config *config) {
		for _, provider := range providers {
			if provider == nil {
				continue
			}
			config.routes = append(config.routes, provider.Routes()...)
		}
	}
}

func muxFromRoutes(routes []Route) http.Handler {
	mux := http.NewServeMux()

	for _, route := range routes {
		if route.handler == nil {
			continue
		}

		mux.Handle(route.pattern, chain(route.handler, route.middlewares...))
	}

	return mux
}
