package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	a "github.com/tony-spark/metrico/internal/agent"
	"github.com/tony-spark/metrico/internal/agent/config"
	"github.com/tony-spark/metrico/internal/agent/metrics"
	"github.com/tony-spark/metrico/internal/agent/transports"
	"github.com/tony-spark/metrico/internal/hash"
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

	var t transports.Transport
	if len(config.Config.Key) > 0 {
		t = transports.NewHTTPTransportHashed(baseURL, hash.NewSha256Hmac(config.Config.Key))
	} else {
		t = transports.NewHTTPTransport(baseURL)
	}

	cs := []metrics.MetricCollector{
		metrics.NewMemoryMetricCollector(),
		metrics.NewRandomMetricCollector(),
		metrics.NewPsUtilMetricsCollector(),
	}

	agent := a.NewMetricsAgent(
		config.Config.PollInterval,
		config.Config.ReportInterval,
		t, cs)
	go agent.Run(context.Background())

	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	log.Info().Msg("Application interrupted via system signal")
}
