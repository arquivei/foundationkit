package app

import (
	"github.com/arquivei/foundationkit/errors"
	"github.com/rs/zerolog/log"
)

func Recover() {
	if r := recover(); r != nil {
		log.Fatal().
			Err(errors.NewFromRecover(r)).
			Msg("[app] App terminated due to panic.")
	}
}
