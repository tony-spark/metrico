package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/server"
	"github.com/tony-spark/metrico/internal/server/config"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	err := config.Parse()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not parse config")
	}

	ctx, cancel := context.WithCancel(context.Background())

	s := server.New(
		server.WithHTTPServer(config.Config.Address),
		server.WithDB(config.Config.DSN),
		server.WithHashKey(config.Config.Key),
		server.WithFileStore(config.Config.StoreFilename, config.Config.StoreInterval, config.Config.Restore),
	)
	go func() {
		err = s.Run(ctx)
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
