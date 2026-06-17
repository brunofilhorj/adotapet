package webserver

import "net/http"

// Middleware wraps an http.Handler with additional behavior.
type Middleware func(http.Handler) http.Handler

// Option customizes the handler built by New.
type Option func(*config)

type config struct {
	handler     http.Handler
	routes      []Route
	middlewares []Middleware
}

func defaultConfig(handler http.Handler) config {
	if handler == nil {
		handler = http.NewServeMux()
	}

	return config{
		handler: handler,
	}
}

// WithMiddlewares applies global middlewares around the root handler.
func WithMiddlewares(middlewares ...Middleware) Option {
	return func(config *config) {
		config.middlewares = append(config.middlewares, middlewares...)
	}
}
