package errors

import (
	"fmt"
)

// KeyValue is used to store a key-value pair within the error
type KeyValue struct {
	Key   interface{}
	Value interface{}
}

func (kv KeyValue) String() string {
	return fmt.Sprintf("%v=%v", kv.Key, kv.Value)
}

// KV is a constructor for KeyValue types
func KV(k, v interface{}) KeyValue {
	return KeyValue{Key: k, Value: v}
}

// GetRootErrorWithKV returns the Err field of Error struct or the error itself if it is of another type
func GetRootErrorWithKV(err error) error {
	var kvs []KeyValue

	for {
		if myErr, ok := err.(Error); ok && myErr.Err != nil {
			if len(myErr.KVs) > 0 {
				kvs = append(kvs, myErr.KVs...)
			}
			err = myErr.Err
			continue
		}
		break
	}
	if len(kvs) == 0 {
		return err
	}
	return E(err, kvs)
}
