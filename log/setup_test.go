package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestMustParseLevel(t *testing.T) {
	assert.Panics(t, func() {
		MustParseLevel("bla")
	})

	for strLevel, zeroLogLevel := range map[string]zerolog.Level{
		"trace": zerolog.TraceLevel,
		"TRACE": zerolog.TraceLevel,
		"Trace": zerolog.TraceLevel,
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
	} {
		assert.Equal(t, zeroLogLevel, MustParseLevel(strLevel))
	}
}

func testLogger(t *testing.T, testFunc func(Config, *bytes.Buffer, string), message string) {
	// prepare data structs
	c := Config{Level: "info"}
	buffer := bytes.Buffer{}
	var parsedLog map[string]*json.RawMessage

	testFunc(c, &buffer, message)

	// parsing log
	_ = json.Unmarshal(buffer.Bytes(), &parsedLog)

	assert.Equal(t, fmt.Sprintf(`"%s"`, message), string(*parsedLog["message"]))
}

func TestSetupLogger(t *testing.T) {
	testLogger(t,
		func(c Config, buffer *bytes.Buffer, msg string) {
			SetupLogger(c, "123", buffer)
			log.Logger.Warn().Msg(msg)
		},
		"bla bla bla")
}

func TestSetupLoggerWithContext(t *testing.T) {
	testLogger(t,
		func(c Config, buffer *bytes.Buffer, msg string) {
			ctx := context.Background()
			ctx = SetupLoggerWithContext(ctx, c, "123", buffer)
			log.Ctx(ctx).Warn().Msg(msg)
		},
		"Bla! Bla! Bla!")
}
