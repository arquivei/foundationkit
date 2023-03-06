package app

import (
	"github.com/rs/zerolog/log"
)

func Recover() {
	if r := recover(); r != nil {
		log.Fatal().
			Err(recoverErr(r)).
			Msg("[app] App terminated due to panic.")
	}
}
