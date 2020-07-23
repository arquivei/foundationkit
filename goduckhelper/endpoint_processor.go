package goduckhelper

import (
	"context"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/goduck"
	"github.com/go-kit/kit/endpoint"
)

// EndpointDecoder decodes a message into a endpoint request.
// See go-kit's endpoint.Endpoint.
type EndpointDecoder func(context.Context, []byte) (interface{}, error)

type processor struct {
	e endpoint.Endpoint
	d EndpointDecoder
}

func (p processor) decodeMessage(ctx context.Context, m []byte) (interface{}, error) {
	const op = errors.Op("decodeMessage")

	r, err := p.d(ctx, m)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return r, nil
}

func (p processor) doEndpoint(ctx context.Context, request interface{}) error {
	const op = errors.Op("doEndpoint")

	_, err := p.e(ctx, request)
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

// Process func will receive the pulled message from the engine.
func (p processor) Process(ctx context.Context, message []byte) error {
	const op = errors.Op("foundationkit/goduckhelper/processor.Process")

	r, err := p.decodeMessage(ctx, message)
	if err != nil {
		return errors.E(op, err)
	}

	err = p.doEndpoint(ctx, r)
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

// NewEndpointProcessor returns a new goduck.Processor that
func NewEndpointProcessor(e endpoint.Endpoint, d EndpointDecoder) (goduck.Processor, error) {
	if e == nil {
		return nil, errors.E("endpoint is nil")
	}

	if d == nil {
		return nil, errors.E("decoder is nil")
	}

	return processor{e: e, d: d}, nil
}

// MustNewEndpointProcessor calls NewEndpointProcessor and panics
// in case of error.
func MustNewEndpointProcessor(e endpoint.Endpoint, d EndpointDecoder) goduck.Processor {
	p, err := NewEndpointProcessor(e, d)
	if err != nil {
		panic(err)
	}
	return p
}
