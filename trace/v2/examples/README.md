# Examples

## API

This example will create an API (ping) that configures the trace, creates a span and propagates to another API (pong). The ping API is configured to use it self as pong API.

To run the example, go to examples folder and run `make run` to create the API. After that, send a request using `make send-without-trace`.
