package loggingmiddleware

import "context"

type requestMetaKeyType int

const requestMetaKey requestMetaKeyType = iota

// GetRequestMeta returns a Meta added in the context by the
// WithRequestMeta function
func GetRequestMeta(ctx context.Context) Meta {
	if v := ctx.Value(requestMetaKey); v != nil {
		return ctx.Value(requestMetaKey).(Meta)
	}
	return nil
}

// WithRequestMeta returns a context the given metadata.
// This value is retrieved by the logging middleware and logged along
// other information.
// This is intended to be used by transport layers to send data to the
// logging middleware on the endpoint layer.
func WithRequestMeta(ctx context.Context, val Meta) context.Context {
	return context.WithValue(ctx, requestMetaKey, val)
}
