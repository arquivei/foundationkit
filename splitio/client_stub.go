package splitio

import "strings"

func parseFeatures(s string) []Feature {
	features := []Feature{}
	for _, feature := range strings.Split(s, ",") {
		features = append(features, Feature(feature))
	}
	return features
}

func mustNewStubClient(features []Feature) Client {
	active := make(map[Feature]struct{}, len(features))
	for _, f := range features {
		active[f] = struct{}{}
	}
	return &stubClient{
		active: active,
	}
}

type stubClient struct {
	active map[Feature]struct{}
}

func (c *stubClient) IsFeatureEnabled(f Feature, _ Attributes) bool {
	_, ok := c.active[f]
	return ok
}

func (c *stubClient) IsFeatureWithUserEnabled(_ User, f Feature, _ Attributes) bool {
	_, ok := c.active[f]
	return ok
}

func (c *stubClient) Close() {
	c.active = nil
}
