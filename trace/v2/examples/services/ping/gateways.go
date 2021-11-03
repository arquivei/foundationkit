package ping

import (
	"context"
	"time"
)

// PongGateway represents the pong gateway
type PongGateway interface {
	Pong(context.Context, int, time.Duration) (string, error)
}
