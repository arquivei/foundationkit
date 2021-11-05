package main

import (
	"time"

	"github.com/arquivei/foundationkit/log"
	"github.com/arquivei/foundationkit/trace/v2"
)

var config struct {
	Log log.Config

	HTTP struct {
		Port string `default:"8686"`
	}
	Shutdown struct {
		GracePeriod time.Duration `default:"3s"`
		Timeout     time.Duration `default:"5s"`
	}
	Pong struct {
		HTTP struct {
			URL     string        `default:"http://localhost:8686/ping/v1"`
			Timeout time.Duration `default:"300s"`
		}
	}

	Trace trace.Config
}
