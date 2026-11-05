package implschemaregistry

import (
	"context"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/schemaregistry"

	"github.com/arquivei/avro/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCacheGetSchemaByID(t *testing.T) {
	mock := new(mockTestRepository)
	mock.On("GetSchemaByID", schemaregistry.ID(10)).
		Return(tagsSchema, nil).
		Times(1)
	mock.On("GetSchemaByID", schemaregistry.ID(99)).
		Return(tagsSchema, errors.New("not found")).
		Times(1)

	repository := WrapWithCache(mock)
	// First time, fetch from repository
	schema, err := repository.GetSchemaByID(context.Background(), 10)
	assert.NoError(t, err)
	assert.Equal(t, tagsSchemaStr, schema.String())

	// Second time, get from schemaByIDCache
	schema, err = repository.GetSchemaByID(context.Background(), 10)
	assert.NoError(t, err)
	assert.Equal(t, tagsSchemaStr, schema.String())

	// Repository returned error
	_, err = repository.GetSchemaByID(context.Background(), 99)
	assert.EqualError(t, errors.GetRootErrorWithKV(err), "not found")

	mock.AssertExpectations(t)
}

func TestCacheGetIDBySchema(t *testing.T) {
	mock := new(mockTestRepository)
	mock.
		On("GetIDBySchema", schemaregistry.Subject("mysubject"), tagsSchemaStr).
		Return(schemaregistry.ID(10), tagsSchema, nil).
		Times(1)
	mock.
		On("GetIDBySchema", schemaregistry.Subject("mysubject"), "wrongschema").
		Return(schemaregistry.ID(10), tagsSchema, errors.New("not found")).
		Times(1)

	repository := WrapWithCache(mock)
	// First time, fetch from repository
	id, schema, err := repository.GetIDBySchema(
		context.Background(),
		schemaregistry.Subject("mysubject"),
		tagsSchemaStr,
	)
	assert.NoError(t, err)
	assert.Equal(t, tagsSchemaStr, schema.String())
	assert.Equal(t, schemaregistry.ID(10), id)

	// Second time, get from idBySchemaCache
	id, schema, err = repository.GetIDBySchema(
		context.Background(),
		schemaregistry.Subject("mysubject"),
		tagsSchemaStr,
	)
	assert.NoError(t, err)
	assert.Equal(t, tagsSchemaStr, schema.String())
	assert.Equal(t, schemaregistry.ID(10), id)

	// Repository returned error
	_, _, err = repository.GetIDBySchema(
		context.Background(),
		schemaregistry.Subject("mysubject"),
		"wrongschema",
	)
	assert.EqualError(t, errors.GetRootErrorWithKV(err), "not found")

	mock.AssertExpectations(t)
}

var cacheLookups = []struct {
	name   string
	lookup func(context.Context, schemaregistry.Repository) error
}{
	{
		name: "GetSchemaByID",
		lookup: func(ctx context.Context, r schemaregistry.Repository) error {
			_, err := r.GetSchemaByID(ctx, 10)

			return err
		},
	},
	{
		name: "GetIDBySchema",
		lookup: func(ctx context.Context, r schemaregistry.Repository) error {
			_, _, err := r.GetIDBySchema(ctx, schemaregistry.Subject("mysubject"), tagsSchemaStr)

			return err
		},
	},
}

func TestCacheConcurrentMissesFetchOnce(t *testing.T) {
	for _, tc := range cacheLookups {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				next := newBlockingTestRepository()
				repository := WrapWithCache(next)

				const callers = 50

				errs := make(chan error, callers)

				for range callers {
					go func() { errs <- tc.lookup(context.Background(), repository) }()
				}

				synctest.Wait()
				close(next.release)

				for range callers {
					assert.NoError(t, <-errs)
				}

				assert.Equal(t, int32(1), next.calls.Load())
			})
		})
	}
}

func TestCacheCanceledCallerDoesNotFailWaiters(t *testing.T) {
	for _, tc := range cacheLookups {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				next := newBlockingTestRepository()
				repository := WrapWithCache(next)

				firstCtx, cancelFirst := context.WithCancel(context.Background())
				firstDone := make(chan struct{})
				go func() {
					_ = tc.lookup(firstCtx, repository)
					close(firstDone)
				}()
				synctest.Wait()

				waiterErr := make(chan error, 1)
				go func() { waiterErr <- tc.lookup(context.Background(), repository) }()
				synctest.Wait()

				cancelFirst()
				synctest.Wait()
				close(next.release)

				<-firstDone

				assert.NoError(t, <-waiterErr)
				assert.Equal(t, int32(1), next.calls.Load())
			})
		})
	}
}

type mockTestRepository struct {
	mock.Mock
}

func (m *mockTestRepository) GetSchemaByID(
	ctx context.Context,
	id schemaregistry.ID,
) (avro.Schema, error) {
	args := m.Called(id)
	return args.Get(0).(avro.Schema), args.Error(1)
}
func (m *mockTestRepository) GetIDBySchema(
	ctx context.Context,
	subject schemaregistry.Subject,
	schema string,
) (schemaregistry.ID, avro.Schema, error) {
	args := m.Called(subject, schema)
	return args.Get(0).(schemaregistry.ID),
		args.Get(1).(avro.Schema), args.Error(2)
}

type blockingTestRepository struct {
	release chan struct{}
	calls   atomic.Int32
}

func newBlockingTestRepository() *blockingTestRepository {
	return &blockingTestRepository{release: make(chan struct{})}
}

func (r *blockingTestRepository) GetSchemaByID(
	ctx context.Context,
	_ schemaregistry.ID,
) (avro.Schema, error) {
	r.calls.Add(1)

	select {
	case <-r.release:
		return tagsSchema, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *blockingTestRepository) GetIDBySchema(
	ctx context.Context,
	_ schemaregistry.Subject,
	_ string,
) (schemaregistry.ID, avro.Schema, error) {
	r.calls.Add(1)

	select {
	case <-r.release:
		return 10, tagsSchema, nil
	case <-ctx.Done():
		return 0, nil, ctx.Err()
	}
}
