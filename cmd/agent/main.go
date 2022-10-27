package main

import (
	"context"
	"flag"
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
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

func main() {
	cfg := config{}

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address to send metrics to")
	flag.DurationVar(&cfg.ReportInterval, "r", 10*time.Second, "report interval")
	flag.DurationVar(&cfg.PollInterval, "p", 2*time.Second, "poll interval")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("Could not parse config")
	}

	log.Printf("Starting metrics agent with config %+v \n", cfg)

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
