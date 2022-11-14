package main

import (
	"context"
	"flag"
	"github.com/caarlos0/env/v6"
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

	cfg := config.Config{}

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address to send metrics to")
	flag.DurationVar(&cfg.ReportInterval, "r", 10*time.Second, "report interval")
	flag.DurationVar(&cfg.PollInterval, "p", 2*time.Second, "poll interval")
	flag.StringVar(&cfg.Key, "k", "", "hash key")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not parse config")
	}

	log.Info().Msgf("Starting metrics agent with config %+v", cfg)

	baseURL := "http://" + strings.Trim(cfg.Address, "\"")

	var t transports.Transport
	if len(cfg.Key) > 0 {
		t = transports.NewHTTPTransportHashed(baseURL, hash.NewSha256Hmac(cfg.Key))
	} else {
		t = transports.NewHTTPTransport(baseURL)
	}

	cs := []metrics.MetricCollector{
		metrics.NewMemoryMetricCollector(),
		metrics.NewRandomMetricCollector(),
		metrics.NewPsUtilMetricsCollector(),
	}

	agent := a.NewMetricsAgent(
		cfg.PollInterval,
		cfg.ReportInterval,
		t, cs)
	go agent.Run(context.Background())

	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	log.Info().Msg("Application interrupted via system signal")
}
