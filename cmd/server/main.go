package main

import (
	"github.com/tony-spark/metrico/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	bindAddress = "127.0.0.1:8080"
)

func main() {
	log.Println("Starting metrics server on", bindAddress)
	go log.Fatal(server.Run(bindAddress))

	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	log.Println("Server interrupted")
}
