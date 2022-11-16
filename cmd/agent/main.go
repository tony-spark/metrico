package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	a "github.com/tony-spark/metrico/internal/agent"
	"github.com/tony-spark/metrico/internal/agent/config"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
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

	go agent.Run(context.Background())

	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	log.Info().Msg("Application interrupted via system signal")
}
