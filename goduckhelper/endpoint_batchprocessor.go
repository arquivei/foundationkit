package goduckhelper

import (
	"context"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/goduck"
	"github.com/go-kit/kit/endpoint"
)

// EndpointBatchDecoder decodes a message into a endpoint request.
// See go-kit's endpoint.Endpoint.
type EndpointBatchDecoder func(context.Context, [][]byte) (interface{}, error)

type batchProcessor struct {
	e endpoint.Endpoint
	d EndpointBatchDecoder
}

func (p batchProcessor) decodeBatch(ctx context.Context, m [][]byte) (interface{}, error) {
	const op = errors.Op("decodeBatch")

	r, err := p.d(ctx, m)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return r, nil
}

func (p batchProcessor) doEndpoint(ctx context.Context, request interface{}) error {
	const op = errors.Op("doEndpoint")

	_, err := p.e(ctx, request)
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

// Process func will receive the pulled message from the engine.
func (p batchProcessor) BatchProcess(ctx context.Context, messages [][]byte) error {
	const op = errors.Op("foundationkit/goduckhelper/batchProcessor.BatchProcess")

	r, err := p.decodeBatch(ctx, messages)
	if err != nil {
		return errors.E(op, err)
	}

	err = p.doEndpoint(ctx, r)
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

// NewEndpointBatchProcessor returns a new goduck.Processor that
func NewEndpointBatchProcessor(e endpoint.Endpoint, d EndpointBatchDecoder) (goduck.BatchProcessor, error) {
	if e == nil {
		return nil, errors.E("endpoint is nil")
	}

	if d == nil {
		return nil, errors.E("decoder is nil")
	}

	return batchProcessor{e: e, d: d}, nil
}

// MustNewEndpointBatchProcessor calls NewEndpointProcessor and panics
// in case of error.
func MustNewEndpointBatchProcessor(e endpoint.Endpoint, d EndpointDecoder) goduck.Processor {
	p, err := NewEndpointProcessor(e, d)
	if err != nil {
		panic(err)
	}
	return p
}
