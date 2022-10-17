package main

import (
	"context"
	a "github.com/tony-spark/metrico/internal/agent"
	t "github.com/tony-spark/metrico/internal/agent/transports"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	baseURL        = "http://127.0.0.1:8080"
)

func main() {
	log.Println("Starting metrics agent...")
	agent := a.NewMetricsAgent(pollInterval, reportInterval, t.NewHTTPTransport(baseURL))
	go agent.Run(context.Background())

	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	log.Println("Application interrupted via system signal")
}
