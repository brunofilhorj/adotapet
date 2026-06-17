package webserver

import "net/http"

// New builds an HTTP handler from routes, services, and global middlewares.
func New(handler http.Handler, options ...Option) http.Handler {
	config := defaultConfig(handler)

	for _, option := range options {
		if option != nil {
			option(&config)
		}
	}

	if len(config.routes) > 0 {
		config.handler = muxFromRoutes(config.routes)
	}

	return chain(config.handler, config.middlewares...)
}

func chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}
