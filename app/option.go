package app

import (
	"net/http"
)

// ServerOption is an interface to Options that can be used to configure a
// http server configured behavior. Ready-to-use Option can be found in
// the `appoptions` package.
type ServerOption func(*http.Server) error

// SetupServer will apply @options on @server, in the order that they
// are declared. If an option returns an error, this function will immediately
// return the received error.
func SetupServer(server *http.Server, options ...ServerOption) error {
	for _, option := range options {
		err := option(server)
		if err != nil {
			return err
		}
	}

	return nil
}
