package app

import (
	"encoding/json"
	"os"

	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/plugins/file"
	"github.com/rs/zerolog/log"
)

// ConfigFilename is the filename of the config file automatically loaded by SetupConfig
var ConfigFilename = "config.json"

// SetupConfig loads the configuration in the given struct. In case of error, prints help and exit application.
func SetupConfig(config interface{}) {
	var files uconfig.Files
	if fileExists(ConfigFilename) {
		files = []struct {
			Path      string
			Unmarshal file.Unmarshal
		}{
			{
				Path:      ConfigFilename,
				Unmarshal: json.Unmarshal,
			},
		}
	}

	c, err := uconfig.Classic(config, files)
	if err != nil {
		if c != nil {
			c.Usage()
		}
		if isErrHelpRequested(err) {
			os.Exit(0)
		}
		log.Fatal().Err(err).Msg("Failed to setup config")
	}
}

func isErrHelpRequested(err error) bool {
	return err != nil && err.Error() == "flag: help requested"
}

func fileExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}
