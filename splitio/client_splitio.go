package splitio

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/arquivei/foundationkit/errors"
	"github.com/splitio/go-client/splitio/client"
	"github.com/splitio/go-client/splitio/conf"
)

type sdkclient struct {
	client *client.SplitClient
	rand   *rand.Rand
}

func (c *sdkclient) IsFeatureEnabled(f Feature, attr Attributes) bool {
	return c.isTreatmentOn(c.getRandomUser(), f, attr)
}

func (c *sdkclient) IsFeatureWithUserEnabled(
	u User, f Feature, attr Attributes) bool {
	return c.isTreatmentOn(u, f, attr)
}

func (c *sdkclient) isTreatmentOn(u User, f Feature, attr Attributes) bool {
	return c.client.Treatment(string(u), string(f), attr) == "on"
}

func (c *sdkclient) Close() {
	c.client.Destroy()
}

func (c *sdkclient) getRandomUser() User {
	user := "user-" + strconv.Itoa(c.rand.Intn(1000))
	return User(user)
}

// mustNewSplitIOClient returns a new Client
func mustNewSplitIOClient(config Config) Client {
	const op = errors.Op("splitio.MustNewSplitIOClient")

	splitConfig := conf.Default()

	splitConfig.LoggerConfig.ErrorWriter = NewZerologLogger()
	splitConfig.LoggerConfig.StandardLoggerFlags = log.Lmsgprefix

	switch strings.ToLower(config.Mode) {
	case "memory":
	case "redis":
		splitConfig.OperationMode = "redis-consumer"
		splitConfig.Redis.Prefix = "FEATURE_FLAG_SPLIT"
		splitConfig.Redis.Database = config.Redis.Database
		splitConfig.Redis.Host = config.Redis.Host
		splitConfig.Redis.Port = config.Redis.Port
	default:
		panic(errors.E(op, "wrong config mode for Split IO"))
	}

	factory, err := client.NewSplitFactory(config.SplitID, splitConfig)
	if err != nil {
		panic(errors.E(op, err))
	}

	splitClient := factory.Client()
	err = splitClient.BlockUntilReady(10)
	if err != nil {
		panic(errors.E(op, err))
	}

	return &sdkclient{
		client: splitClient,
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}
