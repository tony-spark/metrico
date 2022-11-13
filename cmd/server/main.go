package main

import (
	"context"
	"flag"
	"github.com/caarlos0/env/v6"
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

	cfg := config.Config{}

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address to listen")
	flag.DurationVar(&cfg.StoreInterval, "i", 300*time.Second, "store interval")
	flag.StringVar(&cfg.StoreFilename, "f", "/tmp/devops-metrics-db.json", "file to persist metrics")
	flag.BoolVar(&cfg.Restore, "r", true, "whether to load metric from file on start")
	flag.StringVar(&cfg.Key, "k", "", "hash key")
	flag.StringVar(&cfg.DSN, "d", "", "database connection string")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal().Msg("Could not parse env config")
	}

	ctx, cancel := context.WithCancel(context.Background())

	log.Info().Msgf("Starting metrics server with config %+v", cfg)
	go log.Fatal().Msg(server.Run(ctx, cfg).Error())

	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	cancel()
	log.Info().Msg("Server interrupted")
}
