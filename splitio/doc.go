/*
Package splitio provides integration with https://www.split.io/.

# Basics

Using the client directly:

	client := splitio.MustNewClient(config)
	attributes := splitIO.Attributes{
		"age": 37,
	}
	if client.IsFeatureEnabled("MY_FEATURE_FLAG", attributes) {
		// do stuff
	}

# Using the middleware

In the service:

	MY_FF := splitio.Feature("MY_FEATURE_FLAG")
	Features := []splitio.Feature{
		MY_FF,
	}

In the transport layer:

	func GetUserFromRequest(ctx context.Context, request interface{}) map[User]Attributes {
		// Extract user and attributes from the endpoint request
	}

In the main package:

	client := splitio.MustNewClient(config)
	middlewareConfig := splitio.DefaultFFMidlewareConfig()
	middlewareConfig.MultiUserDecodeFn = myapi.GetUserFromRequest
	middlewareConfig.Features = myservice.Features

	middleware := NewFeatureFlagMiddleware(client, middlewareConfig)
	myEndpoint := endpoint.Chain(
		// ...
		middleware,
		// ...
	)(myEndpoint)

In your code:

	if splitio.IsFeatureEnabled(ctx, MY_FF) {
		// do stuff
	}

Although the initial setup is more complex, it has the advantage of
setting up everything only once, and then integrating seamlessly
with new feature flags.

Each feature is checked once per request, for the users and
attributes specified in the "MultiUserDecodeFn". The behavior is
stored in the context, so that it is possible to check anywhere
in the service if the feature is enabled without caring about
which user or attributes should be used.
*/
package splitio
