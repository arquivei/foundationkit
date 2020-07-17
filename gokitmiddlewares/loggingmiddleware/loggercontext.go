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
	// If there is a logger in the context, it's not updated
	return l.WithContext(ctx)
}

func enrichLoggerContext(ctx context.Context, name string, c Config, req interface{}) {
	l := log.Ctx(ctx)
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

		if rid := request.GetRequestIDFromContext(ctx); !rid.IsEmpty() {
			zctx = zctx.EmbedObject(request.GetRequestIDFromContext(ctx))
		} else {
			log.Warn().Msg("Request doesn't have a Request ID! Did you forget to use trackingmiddleware on the trasport layer?")
		}

		if t := trace.GetTraceFromContext(ctx); !trace.IDIsEmpty(t.ID) {
			zctx = zctx.EmbedObject(trace.GetTraceFromContext(ctx))
		} else {
			log.Warn().Msg("Reqeust doesn't have a trace! Did you forget to use trackingmiddleware on the trasport layer?")
		}

		return zctx.Str("endpoint_name", name)

	})
}

func enrichLoggerAfterResponse(ctx context.Context, c Config, begin time.Time, resp interface{}) {
	l := log.Ctx(ctx)

	l.UpdateContext(func(zctx zerolog.Context) zerolog.Context {
		if l.WithLevel(c.LogResponseIfLevel).Enabled() {
			zctx = zctx.Str("endpoint_response", toString(resp, c.TruncResponseAt))
		}
		return zctx.Dur("endpoint_took", time.Since(begin))
	})
}
