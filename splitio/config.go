package splitio

// Config represents the informations to create a Split Client
type Config struct {
	Provider string `default:"splitio"`
	SplitID  string `secret:"true"`
	Mode     string `default:"memory"`
	Redis    struct {
		Database int
		Host     string
		Port     int
	}
	Stub struct {
		Active string `default:""`
	}
}
