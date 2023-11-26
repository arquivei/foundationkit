package auditmiddleware

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/arquivei/foundationkit/message"
	"github.com/rs/zerolog/log"
)

type auditResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newAuditResponseWriter(w http.ResponseWriter) *auditResponseWriter {
	return &auditResponseWriter{w, http.StatusOK}
}

// WriteHeader is used to store the status code that will be returned by the server
func (arw *auditResponseWriter) WriteHeader(code int) {
	arw.statusCode = code
	arw.ResponseWriter.WriteHeader(code)
}

// Exporter is used to push the audit message that was created using CreateAuditMessageFunc
type Exporter interface {
	Push(context.Context, message.Message) error
}

// CreateAuditMessageFunc is used to customize the creation of an audit message
// using http.Request and response status code.
type CreateAuditMessageFunc func(
	request *http.Request,
	responseStatusCode int,
) (message.Message, error)

// AuditMiddleware allows to create and push a customized message for every http request received
type AuditMiddleware struct {
	exporter               Exporter
	createAuditMessageFunc CreateAuditMessageFunc
}

// NewAuditMiddleware returns a new AuditMiddleware with the
// given Exporter and CreateAuditMessageFunc
// Parameters:
//   - @exporter: exporter that will be used to push audit messages.
//   - @createAuditMessageFunc: function that will be used to create audit messages
func NewAuditMiddleware(
	exporter Exporter,
	createAuditMessageFunc CreateAuditMessageFunc,
) *AuditMiddleware {
	return &AuditMiddleware{exporter, createAuditMessageFunc}
}

// New returns a new audit http handler wrapping the @next handler.
func (a *AuditMiddleware) New(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		arw := newAuditResponseWriter(w)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Warn().Err(err).Msg("failed to get request body")
			body = []byte{}
		}
		r.Body.Close()

		requestClone := r.Clone(r.Context())
		r.Body = io.NopCloser(bytes.NewReader(body))
		requestClone.Body = io.NopCloser(bytes.NewReader(body))

		defer func() {
			auditMessage, err := a.createAuditMessageFunc(
				requestClone,
				arw.statusCode,
			)
			if err != nil {
				log.Warn().Err(err).Msg("failed to create audit message")
				return
			}

			err = a.exporter.Push(context.Background(), auditMessage)
			if err != nil {
				log.Warn().Err(err).Msg("failed to send audit message")
				return
			}
		}()

		next.ServeHTTP(arw, r)
	})
}
