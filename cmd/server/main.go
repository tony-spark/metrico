package main

import (
	"context"
	"flag"
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
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFilename string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

func main() {
	cfg := config{}

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address to listen")
	flag.DurationVar(&cfg.StoreInterval, "i", 300*time.Second, "store interval")
	flag.StringVar(&cfg.StoreFilename, "f", "/tmp/devops-metrics-db.json", "file to persist metrics")
	flag.BoolVar(&cfg.Restore, "r", true, "whether to load metric from file on start")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("Could not parse env config")
	}

	ctx, cancel := context.WithCancel(context.Background())

	log.Printf("Starting metrics server with config %+v\n", cfg)
	go log.Fatal(server.Run(ctx, strings.Trim(cfg.Address, "\""), cfg.StoreFilename, cfg.Restore, cfg.StoreInterval))

	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	cancel()
	log.Println("Server interrupted")
}
