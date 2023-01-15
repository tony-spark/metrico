// Package config contains application configuration support (via program arguments and environment variables)
package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

var (
	Config config
)

type config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFilename string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
	DSN           string        `env:"DATABASE_DSN"`
}

func Parse() error {
	flag.StringVar(&Config.Address, "a", "127.0.0.1:8080", "address to listen")
	flag.DurationVar(&Config.StoreInterval, "i", 300*time.Second, "store interval")
	flag.StringVar(&Config.StoreFilename, "f", "/tmp/devops-metrics-db.json", "file to persist metrics")
	flag.BoolVar(&Config.Restore, "r", true, "whether to load metric from file on start")
	flag.StringVar(&Config.Key, "k", "", "hash key")
	flag.StringVar(&Config.DSN, "d", "", "database connection string")
	flag.Parse()

	err := env.Parse(&Config)
	if err != nil {
		return err
	}

	log.Info().Msgf("Server config parsed:  %+v", Config)
	return nil
}
