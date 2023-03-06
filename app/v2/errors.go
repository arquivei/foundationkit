package app

import (
	"fmt"
)

func recoverErr(r any) error {
	if r == nil {
		return nil
	}
	if err, ok := r.(error); ok {
		return fmt.Errorf("panic: %w", err)
	}
	return fmt.Errorf("panic: %v", r)
}
