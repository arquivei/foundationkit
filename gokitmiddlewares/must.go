package gokitmiddlewares

import "github.com/go-kit/kit/endpoint"

// Must returns the endpoint or panics in case of error
//
// This is a helper hor wrapping a New middleware function.
func Must(e endpoint.Middleware, err error) endpoint.Middleware {
	if err != nil {
		panic(err)
	}
	return e
}
