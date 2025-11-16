package main

import (
	"github.com/rs/zerolog/log"

	"github.com/fikriahmadf/outbox-examples/internal/app"
)

func main() {
	application, err := app.InitializeApplication()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize application")
	}

	log.Info().Msg("application initialized")
	_ = application
}
