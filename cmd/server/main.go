package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/tony-spark/metrico/internal/server"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type config struct {
	Address       string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFilename string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func main() {
	cfg := config{}

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("Could not parse config")
	}

	ctx, cancel := context.WithCancel(context.Background())

	log.Println("Starting metrics server on", cfg.Address)
	go log.Fatal(server.Run(ctx, strings.Trim(cfg.Address, "\""), cfg.StoreFilename, cfg.Restore, cfg.StoreInterval))

	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	cancel()
	log.Println("Server interrupted")
}
