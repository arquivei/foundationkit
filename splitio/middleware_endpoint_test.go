package splitio

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClient struct {
	mock.Mock
}

func (c *mockClient) IsFeatureEnabled(f Feature, attr Attributes) bool {
	args := c.Called(f, attr)
	return args.Bool(0)
}

func (c *mockClient) IsFeatureWithUserEnabled(
	u User, f Feature, attr Attributes,
) bool {
	args := c.Called(u, f, attr)
	return args.Bool(0)
}

func (c *mockClient) Close() {
	c.Called()
}

func TestMiddlewareNoUserFn(t *testing.T) {
	tests := []struct {
		name     string
		client   *mockClient
		features []Feature
		checks   map[Feature]bool
	}{
		{
			name: "true",
			client: func() *mockClient {
				client := new(mockClient)
				client.On("IsFeatureEnabled", Feature("TEST"), Attributes{}).Return(true)
				return client
			}(),
			features: []Feature{"TEST"},
			checks: map[Feature]bool{
				Feature("TEST"): true,
			},
		},
		{
			name: "false",
			client: func() *mockClient {
				client := new(mockClient)
				client.On("IsFeatureEnabled", Feature("TEST"), Attributes{}).Return(false)
				return client
			}(),
			features: []Feature{"TEST"},
			checks: map[Feature]bool{
				Feature("TEST"): false,
			},
		},
		{
			name: "multiple",
			client: func() *mockClient {
				client := new(mockClient)
				client.
					On("IsFeatureEnabled", Feature("TEST1"), Attributes{}).Return(false).
					On("IsFeatureEnabled", Feature("TEST2"), Attributes{}).Return(true).
					On("IsFeatureEnabled", Feature("TEST3"), Attributes{}).Return(false).
					On("IsFeatureEnabled", Feature("TEST4"), Attributes{}).Return(true)
				return client
			}(),
			features: []Feature{"TEST1", "TEST2", "TEST3", "TEST4"},
			checks: map[Feature]bool{
				Feature("TEST1"): false,
				Feature("TEST2"): true,
				Feature("TEST3"): false,
				Feature("TEST4"): true,
			},
		},
		{
			name: "feature not checked defaults to false",
			client: func() *mockClient {
				client := new(mockClient)
				client.
					On("IsFeatureEnabled", Feature("TEST1"), Attributes{}).Return(true)
				return client
			}(),
			features: []Feature{"TEST1"},
			checks: map[Feature]bool{
				Feature("TEST1"): true,
				Feature("TEST2"): false,
			},
		},
		{
			name: "empty",
			client: func() *mockClient {
				client := new(mockClient)
				return client
			}(),
			features: []Feature{},
			checks: map[Feature]bool{
				Feature("TEST"): false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			endpoint := func(ctx context.Context, _ interface{}) (interface{}, error) {
				for feature, shouldBeEnabled := range test.checks {
					assert.Equal(t, shouldBeEnabled, IsFeatureEnabled(ctx, feature), feature)
					// can't check trace labels
				}
				return nil, nil
			}
			config := DefaultFFMidlewareConfig()
			config.Features = test.features

			middleware := NewFeatureFlagMiddleware(test.client, config)
			endpoint = middleware(endpoint)
			endpoint(context.Background(), nil)
		})
	}
}

func TestMiddlewareMultiUserFn(t *testing.T) {
	type mockRequest struct {
		users map[User]Attributes
	}
	tests := []struct {
		name     string
		client   *mockClient
		features []Feature
		request  mockRequest
		checks   map[Feature]bool
	}{
		{
			name: "true",
			client: func() *mockClient {
				client := new(mockClient)
				client.On("IsFeatureWithUserEnabled", User("myuser1"), Feature("TEST"), Attributes{"attr": "test1"}).Return(true)
				client.On("IsFeatureWithUserEnabled", User("myuser2"), Feature("TEST"), Attributes{"attr": "test2"}).Return(true)
				return client
			}(),
			features: []Feature{"TEST"},
			request: mockRequest{
				users: map[User]Attributes{
					"myuser1": {"attr": "test1"},
					"myuser2": {"attr": "test2"},
				},
			},
			checks: map[Feature]bool{
				Feature("TEST"): true,
			},
		},
		{
			name: "false",
			client: func() *mockClient {
				client := new(mockClient)
				client.On("IsFeatureWithUserEnabled", User("myuser1"), Feature("TEST"), Attributes{"attr": "test1"}).Return(false)
				client.On("IsFeatureWithUserEnabled", User("myuser2"), Feature("TEST"), Attributes{"attr": "test2"}).Return(false)
				return client
			}(),
			features: []Feature{"TEST"},
			request: mockRequest{
				users: map[User]Attributes{
					"myuser1": {"attr": "test1"},
					"myuser2": {"attr": "test2"},
				},
			},
			checks: map[Feature]bool{
				Feature("TEST"): false,
			},
		},
		{
			name: "multiple",
			client: func() *mockClient {
				client := new(mockClient)
				client.
					On("IsFeatureWithUserEnabled", User("myuser1"), Feature("TEST1"), Attributes{"attr": "test1"}).Return(false).
					On("IsFeatureWithUserEnabled", User("myuser1"), Feature("TEST2"), Attributes{"attr": "test1"}).Return(true).
					On("IsFeatureWithUserEnabled", User("myuser1"), Feature("TEST3"), Attributes{"attr": "test1"}).Return(false).
					On("IsFeatureWithUserEnabled", User("myuser1"), Feature("TEST4"), Attributes{"attr": "test1"}).Return(true).
					On("IsFeatureWithUserEnabled", User("myuser2"), Feature("TEST1"), Attributes{"attr": "test2"}).Return(false).
					On("IsFeatureWithUserEnabled", User("myuser2"), Feature("TEST2"), Attributes{"attr": "test2"}).Return(true).
					On("IsFeatureWithUserEnabled", User("myuser2"), Feature("TEST3"), Attributes{"attr": "test2"}).Return(false).
					On("IsFeatureWithUserEnabled", User("myuser2"), Feature("TEST4"), Attributes{"attr": "test2"}).Return(true)
				return client
			}(),
			features: []Feature{"TEST1", "TEST2", "TEST3", "TEST4"},
			request: mockRequest{
				users: map[User]Attributes{
					"myuser1": {"attr": "test1"},
					"myuser2": {"attr": "test2"},
				},
			},
			checks: map[Feature]bool{
				Feature("TEST1"): false,
				Feature("TEST2"): true,
				Feature("TEST3"): false,
				Feature("TEST4"): true,
			},
		},
		{
			name: "no features",
			client: func() *mockClient {
				client := new(mockClient)
				return client
			}(),
			features: []Feature{},
			request: mockRequest{
				users: map[User]Attributes{
					"myuser": {"attr": "test"},
				},
			},
			checks: map[Feature]bool{
				Feature("TEST"): false,
			},
		},
		{
			name: "no user - defaults to disabled",
			client: func() *mockClient {
				client := new(mockClient)
				return client
			}(),
			features: []Feature{"TEST"},
			request: mockRequest{
				users: map[User]Attributes{},
			},
			checks: map[Feature]bool{
				Feature("TEST"): false,
			},
		},
		{
			name: "multiuser - different behavior",
			client: func() *mockClient {
				client := new(mockClient)
				client.
					On("IsFeatureWithUserEnabled", User("myuser1"), Feature("TEST1"), Attributes{"attr": "test1"}).Return(true).
					On("IsFeatureWithUserEnabled", User("myuser1"), Feature("TEST2"), Attributes{"attr": "test1"}).Return(true).
					On("IsFeatureWithUserEnabled", User("myuser2"), Feature("TEST1"), Attributes{"attr": "test2"}).Return(false).
					On("IsFeatureWithUserEnabled", User("myuser2"), Feature("TEST2"), Attributes{"attr": "test2"}).Return(true)
				return client
			}(),
			features: []Feature{"TEST1", "TEST2"},
			request: mockRequest{
				users: map[User]Attributes{
					"myuser1": {"attr": "test1"},
					"myuser2": {"attr": "test2"},
				},
			},
			checks: map[Feature]bool{
				Feature("TEST1"): false,
				Feature("TEST2"): true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			endpoint := func(ctx context.Context, _ interface{}) (interface{}, error) {
				for feature, shouldBeEnabled := range test.checks {
					assert.Equal(t, shouldBeEnabled, IsFeatureEnabled(ctx, feature), feature)
					// can't check trace labels
				}
				return nil, nil
			}
			config := DefaultFFMidlewareConfig()
			config.Features = test.features
			config.MultiUserDecodeFn = func(_ context.Context, request interface{}) map[User]Attributes {
				r := request.(mockRequest)
				return r.users
			}

			middleware := NewFeatureFlagMiddleware(test.client, config)
			endpoint = middleware(endpoint)
			endpoint(context.Background(), test.request)
		})
	}
}
