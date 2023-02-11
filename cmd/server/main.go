// Package main contains entrypoint for server application
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/tony-spark/metrico/internal/server"
	"github.com/tony-spark/metrico/internal/server/config"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

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
		server.WithCryptoKey(config.Config.PrivateKeyFile),
	)

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	shutdownDone := make(chan struct{})

	go func() {
		<-shutdownSignal

		log.Info().Msg("shutting down gracefully...")
		err = s.Shutdown(context.Background())
		if err != nil {
			log.Fatal().Err(err).Msg("failed to shut down gracefully")
		}

		cancel()
		close(shutdownDone)
	}()

	err = s.Run(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error running server")
	}

	<-shutdownDone
	log.Info().Msg("server shut down gracefully")
}
