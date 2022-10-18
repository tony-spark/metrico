package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/tony-spark/metrico/internal/server"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type config struct {
	Address string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
}

func main() {
	cfg := config{}

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("Could not parse config")
	}

	log.Println("Starting metrics server on", cfg.Address)
	go log.Fatal(server.Run(strings.Trim(cfg.Address, "\"")))

	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	log.Println("Server interrupted")
}
