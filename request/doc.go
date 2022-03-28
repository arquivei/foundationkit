/*
Package request provides helper functions to handle Request ID propagation.

Basics

It will be used the following service as example:

	type Service interface {
		Do(context.Context, Request) (Response, error)
	}
	type Response struct {
	//(...)
		RequestID request.ID
	}

HTTP Layer

Use the function "WithID" to create and put a Request ID in context:

	func MakeEndpoint(s Service) endpoint.Endpoint {
		return func(ctx context.Context, r interface{}) (interface{}, error) {
			req := r.(Request)
			ctx = request.WithID(ctx)
			response, err := s.Do(ctx, req)
			return response, err
		}
	}

Logging

Use the function "GetIDFromContext" to log the Request ID:

	func (l *logging) Do(ctx ontext.Context, req Request) (response Response, err error) {
		logger := log.Logger.With().
			EmbedObject(request.GetIDFromContext(ctx)).
			Logger()
		// (...)
	}

*/
package request
