package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/server"
	"github.com/tony-spark/metrico/internal/server/config"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	err := config.Parse()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not parse config")
	}

	ctx, cancel := context.WithCancel(context.Background())

	log.Info().Msgf("Starting metrics server with config %+v", config.Config)
	go func() {
		err = server.Run(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("Error running server")
		}
	}()

	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	cancel()
	log.Info().Msg("Server interrupted")
}
