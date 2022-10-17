package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	a "github.com/tony-spark/metrico/internal/agent"
	t "github.com/tony-spark/metrico/internal/agent/transports"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	Address        string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	ReportInterval uint   `env:"REPORT_INTERVAL" envDefault:"10"`
	PollInterval   uint   `env:"REPORT_INTERVAL" envDefault:"2"`
}

func main() {
	cfg := config{}

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("Could not parse config")
	}

	log.Printf("Starting metrics agent %+v \n", cfg)

	agent := a.NewMetricsAgent(time.Duration(cfg.PollInterval)*time.Second, time.Duration(cfg.ReportInterval)*time.Second, t.NewHTTPTransport(cfg.Address))
	go agent.Run(context.Background())

	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	log.Println("Application interrupted via system signal")
}
