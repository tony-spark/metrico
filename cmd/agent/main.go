// Package main contains entrypoint for agent application
package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/agent/transports"
	"github.com/tony-spark/metrico/internal/hash"

	a "github.com/tony-spark/metrico/internal/agent"
	"github.com/tony-spark/metrico/internal/agent/config"
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

	baseURL := "http://" + strings.Trim(config.Config.Address, "\"")

	var transportOptions []transports.HTTPOption
	if len(config.Config.Key) > 0 {
		transportOptions = append(transportOptions, transports.WithHasher(hash.NewSha256Hmac(config.Config.Key)))
	}
	t := transports.NewHTTP(baseURL, transportOptions...)

	agent := a.New(
		a.WithTransport(t),
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
