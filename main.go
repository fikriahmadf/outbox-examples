package main

import (
	"github.com/rs/zerolog/log"

	"github.com/fikriahmadf/outbox-examples/internal/app"
)

func main() {
	server, err := app.InitializeServer()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize server")
	}
	log.Info().Msg("server initialized")
	server.SetupAndServe()
}
