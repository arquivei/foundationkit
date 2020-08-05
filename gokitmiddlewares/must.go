package gokitmiddlewares

import "github.com/go-kit/kit/endpoint"

func Must(e endpoint.Middleware, err error) endpoint.Middleware {
	if err != nil {
		panic(err)
	}
	return e
}
