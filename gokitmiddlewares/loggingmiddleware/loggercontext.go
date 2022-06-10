package loggingmiddleware

import (
	"context"
	"time"

	"github.com/arquivei/foundationkit/request"
	"github.com/arquivei/foundationkit/trace"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initLoggerContext(ctx context.Context, l zerolog.Logger) context.Context {
	// Creates a copy, otherwise the logger is updated for every request
	// Ensures that there is a logger in the context
	// If the same logger is in the context already or if there is a disabled
	// logger already in the context, the context is not updated.
	logger := l.With().Logger()
	return logger.WithContext(ctx)
}

func enrichLoggerContext(ctx context.Context, l *zerolog.Logger, c Config, req interface{}) {
	l.UpdateContext(func(zctx zerolog.Context) zerolog.Context {
		if c.Meta != nil {
			zctx = zctx.Interface("endpoint_meta", c.Meta)
		}

		if l.WithLevel(c.LogRequestIfLevel).Enabled() {
			zctx = zctx.Str("endpoint_request", toString(req, c.TruncRequestAt))
		}

		if meta := GetRequestMeta(ctx); len(meta) > 0 {
			zctx = zctx.Interface("request_meta", meta)
		}

		if rid := request.GetIDFromContext(ctx); !rid.IsEmpty() {
			zctx = zctx.EmbedObject(request.GetIDFromContext(ctx))
		} else {
			log.Warn().Msg("Request doesn't have a Request ID! Did you forget to use trackingmiddleware on the transport layer?")
		}

		if t := trace.GetFromContext(ctx); !trace.IDIsEmpty(t.ID) {
			zctx = zctx.EmbedObject(trace.GetFromContext(ctx))
		} else {
			log.Warn().Msg("Request doesn't have a trace! Did you forget to use trackingmiddleware on the transport layer?")
		}

		return zctx.Str("endpoint_name", c.Name)
	})
}

func enrichLoggerAfterResponse(l *zerolog.Logger, c Config, begin time.Time, resp interface{}) {
	l.UpdateContext(func(zctx zerolog.Context) zerolog.Context {
		if l.WithLevel(c.LogResponseIfLevel).Enabled() {
			zctx = zctx.Str("endpoint_response", toString(resp, c.TruncResponseAt))
		}
		return zctx.Dur("endpoint_took", time.Since(begin))
	})
}
