package implping

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/trace/v2"
	"github.com/arquivei/foundationkit/trace/v2/examples/services/ping"
)

type adapterPongHTTP struct {
	client *http.Client
	url    string
}

type pongRequestResponse struct {
	Num   int           `json:"num"`
	Sleep time.Duration `json:"sleep"`
}

// NewHTTPPongAdapter returns pong gateway implementation using http
func NewHTTPPongAdapter(
	client *http.Client,
	url string,
) ping.PongGateway {
	return &adapterPongHTTP{
		client: client,
		url:    url,
	}
}

func (a *adapterPongHTTP) Pong(
	ctx context.Context,
	num int,
	sleep time.Duration,
) (string, error) {
	const op = errors.Op("implping.adapterPongHttp.Pong")
	ctx, span := trace.Start(ctx, "implping.adapterPongHttp.Pong")
	defer span.End()

	body, err := json.Marshal(pongRequestResponse{
		Num:   num,
		Sleep: sleep,
	})
	if err != nil {
		return "", errors.E(err, op)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		a.url,
		bytes.NewReader(body),
	)
	if err != nil {
		return "", errors.E(err, op)
	}

	trace.SetTraceInRequest(request)

	httpResponse, err := a.client.Do(request)
	if err != nil {
		return "", errors.E(err, op)
	}
	defer httpResponse.Body.Close()

	var response string
	return response, json.NewDecoder(httpResponse.Body).Decode(&response)
}
