package timeoutmiddleware

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arquivei/foundationkit/errors"
)

func newEndpoint(t *testing.T, c Config, next endpoint.Endpoint) endpoint.Endpoint {
	t.Helper()
	middleware, err := New(c)
	require.NoError(t, err)
	return middleware(next)
}

func TestTimeoutMiddleware(t *testing.T) {
	someError := errors.E("some endpoint error")

	tests := []struct {
		name             string
		config           Config
		next             endpoint.Endpoint
		expectedResponse interface{}
		expectedError    error
		expectsDeadline  bool
	}{
		{
			name:   "endpoint responds within the deadline",
			config: NewDefaultConfig(),
			next: func(ctx context.Context, request interface{}) (interface{}, error) {
				return "ok", nil
			},
			expectedResponse: "ok",
		},
		{
			name:   "endpoint error within the deadline is returned unchanged",
			config: NewDefaultConfig(),
			next: func(ctx context.Context, request interface{}) (interface{}, error) {
				return nil, someError
			},
			expectedError: someError,
		},
		{
			name: "slow endpoint with Wait false returns deadline error",
			config: Config{
				Timeout:       20 * time.Millisecond,
				Wait:          false,
				ErrorSeverity: errors.SeverityRuntime,
				ErrorCode:     errors.Code("TIMEOUT"),
			},
			next: func(ctx context.Context, request interface{}) (interface{}, error) {
				<-ctx.Done()
				time.Sleep(50 * time.Millisecond)
				return "too late", nil
			},
			expectsDeadline: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := newEndpoint(t, test.config, test.next)(context.Background(), nil)

			if test.expectsDeadline {
				require.Error(t, err)
				assert.ErrorIs(t, err, context.DeadlineExceeded)
				assert.Equal(t, test.config.ErrorSeverity, errors.GetSeverity(err))
				assert.Equal(t, test.config.ErrorCode, errors.GetCode(err))
				assert.Nil(t, response)
				return
			}

			assert.Equal(t, test.expectedResponse, response)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

// TestTimeoutMiddlewareDoesNotLeakGoroutines is a regression test: with an
// unbuffered result channel, every timed-out request left the inner
// go-routine blocked forever on the channel send, holding the endpoint
// response in memory until the process died.
func TestTimeoutMiddlewareDoesNotLeakGoroutines(t *testing.T) {
	tests := []struct {
		name string
		next func(release <-chan struct{}) endpoint.Endpoint
	}{
		{
			name: "endpoint finishes after the timeout",
			next: func(release <-chan struct{}) endpoint.Endpoint {
				return func(ctx context.Context, request interface{}) (interface{}, error) {
					<-release
					return "late response", nil
				}
			},
		},
		{
			name: "endpoint panics after the timeout",
			next: func(release <-chan struct{}) endpoint.Endpoint {
				return func(ctx context.Context, request interface{}) (interface{}, error) {
					<-release
					panic("panic after timeout")
				}
			},
		},
	}

	config := Config{
		Timeout:       10 * time.Millisecond,
		Wait:          false,
		ErrorSeverity: errors.SeverityRuntime,
		ErrorCode:     errors.Code("TIMEOUT"),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			baseline := runtime.NumGoroutine()

			release := make(chan struct{})
			_, err := newEndpoint(t, config, test.next(release))(context.Background(), nil)
			require.Error(t, err)
			assert.ErrorIs(t, err, context.DeadlineExceeded)

			// The caller already got the timeout error. Once the endpoint
			// finishes, the inner go-routine must terminate as well.
			// The polling runs inline (not testify's Eventually) because
			// Eventually evaluates the condition in its own go-routine,
			// which would inflate runtime.NumGoroutine by one.
			close(release)
			assert.LessOrEqual(
				t,
				waitNumGoroutineAtMost(baseline, 2*time.Second),
				baseline,
				"inner go-routine is still alive after the endpoint finished",
			)
		})
	}
}

// waitNumGoroutineAtMost polls until the number of go-routines drops to max
// or the timeout expires, and returns the last observed count.
func waitNumGoroutineAtMost(max int, timeout time.Duration) int {
	deadline := time.Now().Add(timeout)
	count := runtime.NumGoroutine()
	for count > max && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
		count = runtime.NumGoroutine()
	}
	return count
}
