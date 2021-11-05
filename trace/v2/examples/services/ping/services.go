package ping

import (
	"context"
	"time"

	"github.com/arquivei/foundationkit/trace/v2"

	"github.com/arquivei/foundationkit/errors"
	"github.com/rs/zerolog/log"
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

	ctx, span := trace.Start(ctx, "ping-service")
	defer span.End()

	t := trace.GetTraceInfoFromContext(ctx)
	defer log.Ctx(ctx).Info().EmbedObject(t).Msg("Just for check trace info")

	time.Sleep(req.Sleep * time.Millisecond)

	if req.Num == 0 {
		return "ping", nil
	}

	pong, err := s.pongGateway.Pong(ctx, req.Num-1, req.Sleep)
	if err != nil {
		return "", errors.E(op, err)
	}

	return "ping " + pong, nil
}
