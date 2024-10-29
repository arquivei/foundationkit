package trackingmiddleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/arquivei/foundationkit/ref"
	"github.com/arquivei/foundationkit/request"
	"github.com/arquivei/foundationkit/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	otrace "go.opentelemetry.io/otel/trace"
)

// New instantiates a new tracking middleware wrapping the @next handler.
func New(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = request.WithNewID(ctx)

		// KLUDGE: This enables migrating from trace v1 to trace v2
		// without breaking compatibility.
		if hasValidTraceV2Header(r) {
			t := translateTraceV2ToTraceV1(r)
			ctx = trace.WithTrace(ctx, t)
		} else {
			// There is no trace v2 in header, let's try with trace v1
			t := trace.GetFromHTTPRequest(r)
			ctx = trace.WithTrace(ctx, t)

			// We fetch trace id from context because WithTrace
			// will initialize a trace if it is empty.
			t = trace.GetFromContext(ctx)
			translateTraceV1ToTraceV2Headers(t, r)
		}

		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(w.Header()))

		request.SetInHTTPResponse(request.GetIDFromContext(ctx), w)
		trace.SetInHTTPResponse(trace.GetFromContext(ctx), w)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func translateTraceV2ToTraceV1(r *http.Request) trace.Trace {
	header := r.Header.Get("traceparent")
	traceInfo := getTraceInfoFromTraceV2Header(header)

	ps, _ := strconv.ParseFloat(traceInfo.probabilitySample, 64)

	return trace.Trace{
		ID:                trace.Parse(traceInfo.traceID),
		ProbabilitySample: ref.Of(ps),
	}
}

func translateTraceV1ToTraceV2Headers(tv1 trace.Trace, r *http.Request) {
	if trace.IDIsEmpty(tv1.ID) {
		return
	}

	tv2, err := otrace.TraceIDFromHex(tv1.ID.String())
	if err != nil {
		return
	}

	// Because we don't have a valid span id, lets fake one using the
	// beginning of the trace id.
	sp := otrace.SpanID(tv2[0:16])

	// For now, the probability is being handled as a boolean. Anything
	// higher than zero will be sampled.
	p := "00"
	if tv1.ProbabilitySample != nil && *tv1.ProbabilitySample == 1 {
		p = "01"
	}
	r.Header.Set("traceparent", "00-"+tv2.String()+"-"+sp.String()+"-"+p)
}

func hasValidTraceV2Header(r *http.Request) bool {
	header := r.Header.Get("traceparent")
	traceInfo := getTraceInfoFromTraceV2Header(header)
	return traceInfo.traceID != ""
}

type traceInfo struct {
	traceID           string
	spanID            string
	probabilitySample string
}

func getTraceInfoFromTraceV2Header(header string) traceInfo {
	// 00-TRACEID-SPANID-PROBABILITYSAMPLE
	s := strings.Split(header, "-")
	if len(s) != 4 {
		return traceInfo{}
	}

	return traceInfo{
		traceID:           s[1],
		spanID:            s[2],
		probabilitySample: s[3],
	}
}
