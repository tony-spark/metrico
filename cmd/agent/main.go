package main

import (
	"github.com/tony-spark/metrico/internal/agent"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	serverAddress  = "http://127.0.0.1:8080"
)

func main() {
	log.Println("Starting metrics agent...")
	go agent.Run(pollInterval, reportInterval, serverAddress)

	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	log.Println("Agent interrupted")
}
