package ping

import (
	"context"
	"fmt"
	"time"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/trace/v2"
	"go.opentelemetry.io/otel/attribute"
)

// Service Service
type Service interface {
	Ping(context.Context, Request) (string, error)
}

type service struct {
	pongGateway PongGateway
}

// NewService NewService
func NewService(pongGateway PongGateway) Service {
	return &service{
		pongGateway: pongGateway,
	}
}

func (s *service) Ping(ctx context.Context, req Request) (string, error) {
	const op = errors.Op("ping.service.Ping")

	ctx, span := trace.Start(ctx, "ping.service.Ping")
	defer span.End()

	if req.Num < 0 {
		return "", fmt.Errorf("negative number received: %d", req.Num)
	}

	time.Sleep(req.Sleep * time.Millisecond)

	pingpong := getPingPong(req.Num)

	span.SetAttributes(
		attribute.KeyValue{Key: "req.num", Value: attribute.IntValue(req.Num)},
		attribute.KeyValue{Key: "ping.pong", Value: attribute.StringValue(pingpong)},
	)

	if req.Num == 0 {
		return pingpong, nil
	}

	pong, err := s.pongGateway.Pong(ctx, req.Num-1, req.Sleep)
	if err != nil {
		return "", errors.E(err, op)
	}

	return pingpong + "-" + pong, nil
}

func getPingPong(n int) string {
	if n%2 == 1 {
		return "ping"
	}
	return "pong"
}
