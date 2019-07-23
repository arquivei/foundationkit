package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVAppending(t *testing.T) {
	err := E(SeverityFatal, "1st", KV("k1", "v1"), KV("k12", "v12"))
	err = E(err, Op("2nd"))
	err = E(Op("3rd"), []KeyValue{KV("k3", "v3")}, err)
	err = E(err, Op("4th"))

	assert.Equal(t, "4th: 3rd: 2nd: 1st [k3=v3,k1=v1,k12=v12]", err.Error())
}

func TestGetRootErrorWithKeyValue(t *testing.T) {
	err := E("a")
	err = E(Op("b"), err)
	err = E(Op("c"), err)
	err = E(Op("d"), KV("k", "v"), err)
	assert.EqualError(t, GetRootErrorWithKV(err), "a [k=v]")
}
