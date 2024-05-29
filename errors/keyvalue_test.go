package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVAppending(t *testing.T) {
	err := New("1st", SeverityFatal, KV("k1", "v1"), KV("k12", "v12"))
	err = E(err, Op("2nd"))
	err = E(err, Op("3rd"), KeyValues{KV("k3", "v3")})
	err = E(err, Op("4th"))

	assert.Equal(t, "4th: 3rd: 2nd: 1st [k3=v3,k1=v1,k12=v12]", err.Error())
}

func TestGetRootErrorWithKeyValue(t *testing.T) {
	err := New("a")
	err = E(err, Op("b"))
	err = E(err, Op("c"))
	err = E(err, Op("d"), KV("k", "v"))
	assert.EqualError(t, GetRootErrorWithKV(err), "a [k=v]")
}
