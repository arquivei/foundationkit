package contextmap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.IsType(t, contextMap{}, New())
}

func Test_contextMap_String(t *testing.T) {
	tests := []struct {
		name string
		m    contextMap
		want string
	}{
		{
			name: "empty map",
			m:    contextMap{},
			want: "{}",
		},
		{
			name: "one string entry",
			m:    contextMap{"key": "value"},
			want: `{"key":"value"}`,
		},
		{
			name: "one int entry",
			m:    contextMap{"key": 1},
			want: `{"key":1}`,
		},
		{
			name: "two entries",
			m:    contextMap{"key1": 1, "key2": 2},
			want: `{"key1":1,"key2":2}`,
		},
		{
			name: "three entries",
			m:    contextMap{"key1": 1, "key2": 2, "key3": "3"},
			want: `{"key1":1,"key2":2,"key3":"3"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.m.String())
		})
	}
}

func Test_contextMap_WithCtx(t *testing.T) {
	m := New()
	ctx := m.WithCtx(context.Background())

	assert.Equal(t, m, Ctx(ctx))
}

func Test_contextMap_SetAndGet(t *testing.T) {
	m := New()

	assert.Nil(t, m.Get("test"), "If value is not on contextMap, returns nil")

	myValue := "pudim"

	m.Set("myvalue", myValue)
	assert.Equal(t, myValue, m.Get("myvalue"), "Get should return the same value set")
}

func TestCtx(t *testing.T) {
	assert.Equal(t, noop, Ctx(context.Background()), "Empty context should return noop")

	m := New()
	ctx := m.WithCtx(context.Background())
	assert.Equal(t, m, Ctx(ctx), "Context with ContextMap")
}
