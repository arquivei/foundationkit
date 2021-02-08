package splitio

import (
	"context"
	"strconv"

	"github.com/arquivei/foundationkit/contextmap"
	"github.com/arquivei/foundationkit/trace"
	"github.com/go-kit/kit/endpoint"
	"github.com/rs/zerolog"
)

// MultiUserDecoder is a function that is able to extract multiple users from
// an endpoint's request, with each user being associated with a set of
// attributes.
// If no user can be extracted from the request, the fuction should return a
// fixed user called "nouser". It's an error to return an empty map, which
// will cause the features to be evaluated as "off"
type MultiUserDecoder func(ctx context.Context, request interface{}) map[User]Attributes

// FFMidlewareConfig contains the parameters for the feature flag middleware
type FFMidlewareConfig struct {
	// MultiUserDecodeFn is a function that extract a multiple users from the
	// endpoint request. The feature flag will be checked against all of the
	// users, and will only be active if all of the users are in the split.
	MultiUserDecodeFn MultiUserDecoder
	// Features are the list of FF that should be checked by the middleware
	Features []Feature
}

// DefaultFFMidlewareConfig returns the default config for the FF Midleware
func DefaultFFMidlewareConfig() FFMidlewareConfig {
	return FFMidlewareConfig{}
}

type ffMidleware struct {
	next   endpoint.Endpoint
	client Client
	config FFMidlewareConfig
}

// NoUser will always have the same behavior
var NoUser = User("nouser")

// NewFeatureFlagMiddleware returns a middleware that checks if a feature is
// enabled.
// The chosen behaviour is stored in the context, and it is possible to check
// its value with splitio.IsFeatureEnabled(ctx, feature)
// It integrates with foudationkit's trace, and stores the behaviour in the
// labels.
func NewFeatureFlagMiddleware(client Client, config FFMidlewareConfig) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		middleware := ffMidleware{
			next:   next,
			client: client,
			config: config,
		}
		return middleware.Next
	}
}

func (m ffMidleware) Next(ctx context.Context, request interface{}) (response interface{}, err error) {
	// Nothing to do :(
	if len(m.config.Features) == 0 {
		return m.next(ctx, request)
	}
	var users map[User]Attributes
	if m.config.MultiUserDecodeFn != nil {
		users = m.config.MultiUserDecodeFn(ctx, request)
	}

	traceLabels := make(map[string]string, len(m.config.Features))
	enabledFeatures := make([]string, 0, len(m.config.Features))

	for _, feature := range m.config.Features {
		isEnabled := m.isEnabled(users, feature)
		traceLabels["feature_"+string(feature)] = strconv.FormatBool(isEnabled)
		ctx = context.WithValue(ctx, feature, isEnabled)
		if isEnabled {
			enabledFeatures = append(enabledFeatures, string(feature))
		}
	}
	ctx = trace.WithLabels(ctx, traceLabels)

	if len(enabledFeatures) > 0 {
		contextmap.Ctx(ctx).Set("enabled_features", enabledFeatures)
	}

	return m.next(ctx, request)
}

func (m *ffMidleware) isEnabled(users map[User]Attributes, feature Feature) bool {
	if m.config.MultiUserDecodeFn != nil {
		return m.isEnabledForUsers(users, feature)
	}
	return m.client.IsFeatureEnabled(feature, Attributes{})
}

func (m *ffMidleware) isEnabledForUsers(users map[User]Attributes, feature Feature) bool {
	if len(users) == 0 {
		return false
	}

	for user, attributes := range users {
		if !m.client.IsFeatureWithUserEnabled(user, feature, attributes) {
			return false
		}
	}
	return true
}

// IsFeatureEnabled checks if the @feature stored in the context @ctx was
// previously selected in the FeatureFlagMiddleware
func IsFeatureEnabled(ctx context.Context, feature Feature) bool {
	value := ctx.Value(feature)
	isEnabled, _ := value.(bool)
	return isEnabled
}

// EnrichLogWithEnabledFeatures takes a zerolog context and adds the key enabled_features with a list
// with all enabled features. If there isn't any feature enabled, it returns the original context.
func EnrichLogWithEnabledFeatures(ctx context.Context, zctx zerolog.Context) zerolog.Context {
	enabledFeatures, _ := contextmap.Ctx(ctx).Get("enabled_features").([]string)
	// We use _ instead of ok because we don't want a panic and also because
	// we use len(enabledFeatures) to test if we should enrich the log.
	if len(enabledFeatures) == 0 {
		return zctx
	}
	return zctx.Strs("enabled_features", enabledFeatures)
}
