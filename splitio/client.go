package splitio

import (
	"strings"

	zlog "github.com/rs/zerolog/log"
)

// Client wraps a splitio client
type Client interface {
	IsFeatureEnabled(Feature, Attributes) bool
	IsFeatureWithUserEnabled(User, Feature, Attributes) bool
	Close()
}

// MustNewClient returns a new Client based on the config
func MustNewClient(config Config) Client {
	switch provider := strings.ToLower(config.Provider); provider {
	case "splitio":
		return mustNewSplitIOClient(config)
	case "stub":
		zlog.Warn().Msg("Using split.io stub")
		return mustNewStubClient(parseFeatures(config.Stub.Active))
	default:
		panic("invalid splitio provider: " + provider)
	}
}
