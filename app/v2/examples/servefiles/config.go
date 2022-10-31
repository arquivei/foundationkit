package main

import "github.com/arquivei/foundationkit/app/v2"

type config struct {
	// App is the app scpecific configuration
	app.Config

	// Programs can have any configuration the want.

	HTTP struct {
		Port string `default:"8000"`
	}
	Dir string `default:"."`
}
