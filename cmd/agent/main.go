package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	a "github.com/tony-spark/metrico/internal/agent"
	t "github.com/tony-spark/metrico/internal/agent/transports"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type config struct {
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

func main() {
	cfg := config{}

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("Could not parse config")
	}

	log.Printf("Starting metrics agent %+v \n", cfg)

	baseURL := "http://" + strings.Trim(cfg.Address, "\"")

	agent := a.NewMetricsAgent(
		cfg.PollInterval,
		cfg.ReportInterval,
		t.NewHTTPTransport(baseURL))
	go agent.Run(context.Background())

	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	log.Println("Application interrupted via system signal")
}
