package contextmap

import (
	"context"
	"encoding/json"
)

type contextmapKeyType int

const contextmapKey contextmapKeyType = 1

// ContextMap provides a way to enrich a context.Context with information on lower layers
// and have this information available on the upper layers.
type ContextMap interface {
	String() string
	WithCtx(context.Context) context.Context
	Set(key string, val interface{}) ContextMap
	Get(key string) interface{}
}

type contextMap map[string]interface{}

// New returns a new ContextMap that can be embeded in a context.Context.
func New() ContextMap {
	return make(contextMap)
}

func (m contextMap) String() string {
	s, err := json.Marshal(m)
	if err != nil {
		return err.Error()
	}
	return string(s)
}

func (m contextMap) WithCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextmapKey, m)
}

func (m contextMap) Set(key string, val interface{}) ContextMap {
	m[key] = val
	return m
}

func (m contextMap) Get(key string) interface{} {
	return m[key]
}

// Ctx retrieved a ContextMap from the context.Context.
// If there isn't a ContextMap in the context, a noop ContextMap is returned, so
// it is safe to do something like Ctx(ctx).Set("key", "value").
func Ctx(ctx context.Context) ContextMap {
	v := ctx.Value(contextmapKey)
	if v == nil {
		return noop
	}
	return v.(ContextMap)
}
