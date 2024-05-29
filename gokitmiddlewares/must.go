package gokitmiddlewares

import "github.com/arquivei/foundationkit/endpoint"

// Must returns the endpoint or panics in case of error
//
// This is a helper hor wrapping a New middleware function.
func Must[Request any, Response any](e endpoint.Middleware[Request, Response], err error) endpoint.Middleware[Request, Response] {
	if err != nil {
		panic(err)
	}
	return e
}
