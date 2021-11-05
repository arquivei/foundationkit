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

type adapterPongHttp struct {
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

	return &adapterPongHttp{
		client: client,
		url:    url,
	}
}

func (a *adapterPongHttp) Pong(
	ctx context.Context,
	num int,
	sleep time.Duration,
) (string, error) {
	const op = errors.Op("implping.adapterPongHttp.Pong")

	body, err := json.Marshal(pongRequestResponse{
		Num:   num,
		Sleep: sleep,
	})
	if err != nil {
		return "", errors.E(op, err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		a.url,
		bytes.NewReader(body),
	)
	if err != nil {
		return "", errors.E(op, err)
	}

	trace.SetTraceInRequest(request)

	response, err := a.client.Do(request)
	if err != nil {
		return "", errors.E(op, err)
	}
	defer response.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)

	return buf.String(), nil
}
