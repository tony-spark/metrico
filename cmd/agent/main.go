package main

import (
	"context"
	"flag"
	"github.com/caarlos0/env/v6"
	a "github.com/tony-spark/metrico/internal/agent"
	"github.com/tony-spark/metrico/internal/agent/config"
	"github.com/tony-spark/metrico/internal/agent/transports"
	"github.com/tony-spark/metrico/internal/hash"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	cfg := config.Config{}

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address to send metrics to")
	flag.DurationVar(&cfg.ReportInterval, "r", 10*time.Second, "report interval")
	flag.DurationVar(&cfg.PollInterval, "p", 2*time.Second, "poll interval")
	flag.StringVar(&cfg.Key, "k", "", "hash key")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("Could not parse config")
	}

	log.Printf("Starting metrics agent with config %+v \n", cfg)

	baseURL := "http://" + strings.Trim(cfg.Address, "\"")

	var t transports.Transport
	if len(cfg.Key) > 0 {
		t = transports.NewHTTPTransportHashed(baseURL, hash.NewSha256Keyed(cfg.Key))
	} else {
		t = transports.NewHTTPTransport(baseURL)
	}

	agent := a.NewMetricsAgent(
		cfg.PollInterval,
		cfg.ReportInterval,
		t)
	go agent.Run(context.Background())

	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	log.Println("Application interrupted via system signal")
}
