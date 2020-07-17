package enrichloggingmiddleware

import (
	"net/http"

	"github.com/arquivei/foundationkit/gokitmiddlewares/loggingmiddleware"
)

func New(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		meta := loggingmiddleware.Meta{
			"url":         r.URL.String(),
			"user_agent":  r.UserAgent(),
			"remote_addr": r.RemoteAddr,
		}

		ctx := r.Context()
		ctx = loggingmiddleware.WithRequestMeta(ctx, meta)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
