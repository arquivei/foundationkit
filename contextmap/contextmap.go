package contextmap

import (
	"context"
	"encoding/json"
	"sync"
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

type contextMap struct {
	mu *sync.RWMutex
	m  map[string]interface{}
}

// New returns a new ContextMap that can be embedded in a context.Context.
func New() ContextMap {
	return contextMap{
		mu: &sync.RWMutex{},
		m:  make(map[string]interface{}),
	}
}

func (cm contextMap) String() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	s, err := json.Marshal(cm.m)
	if err != nil {
		return err.Error()
	}
	return string(s)
}

func (cm contextMap) WithCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextmapKey, cm)
}

func (cm contextMap) Set(key string, val interface{}) ContextMap {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.m[key] = val
	return cm
}

func (cm contextMap) Get(key string) interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.m[key]
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
