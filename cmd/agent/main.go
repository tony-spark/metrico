// Package main contains entrypoint for agent application
package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	a "github.com/tony-spark/metrico/internal/agent"
	"github.com/tony-spark/metrico/internal/agent/config"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	err := config.Parse()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not parse config")
	}

	baseURL := "http://" + strings.Trim(config.Config.Address, "\"")

	agent := a.New(
		a.WithHTTPTransport(baseURL, config.Config.Key),
		a.WithPollInterval(config.Config.PollInterval),
		a.WithReportInterval(config.Config.ReportInterval),
	)

	ctx, cancel := context.WithCancel(context.Background())
	go agent.Run(ctx)

	if config.Config.Profile {
		log.Info().Msg("starting profile http server")
		go func() {
			err := http.ListenAndServe("127.0.0.1:8888", nil)
			if err != nil {
				log.Error().Err(err).Msg("error running http server")
			}
		}()
	}

	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	cancel()
	log.Info().Msg("Application interrupted via system signal")
}
