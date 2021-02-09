package contextmap

import "context"

var noop noopContextMap

type noopContextMap struct{}

func (m noopContextMap) String() string {
	return ""
}

func (m noopContextMap) WithCtx(ctx context.Context) context.Context {
	return ctx
}

func (m noopContextMap) Set(key string, val interface{}) ContextMap {
	return m
}

func (m noopContextMap) Get(key string) interface{} {
	return nil
}
