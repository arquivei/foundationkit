package appoptions

import (
	"net/http"
	"time"

	"github.com/arquivei/foundationkit/app"
)

// WithReadTimeout will apply @timeout to ReadTimeout of the http server
func WithReadTimeout(timeout time.Duration) app.ServerOption {
	return func(server *http.Server) error {
		server.ReadTimeout = timeout
		return nil
	}
}

// WithReadHeaderTimeout will apply @timeout to ReadHeaderTimeout of the
// http server
func WithReadHeaderTimeout(timeout time.Duration) app.ServerOption {
	return func(server *http.Server) error {
		server.ReadHeaderTimeout = timeout
		return nil
	}
}

// WithReadTimeout will apply @timeout to WriteTimeout of the http server
func WithWriteTimeout(timeout time.Duration) app.ServerOption {
	return func(server *http.Server) error {
		server.WriteTimeout = timeout
		return nil
	}
}
