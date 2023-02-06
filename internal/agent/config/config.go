// Package config contains agent application configuration support (via program arguments and environment variables)
package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

var (
	Config config
)

type config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
	Profile        bool          `env:"PROFILING"`
	PublicKeyFile  string        `env:"CRYPTO_KEY"`
}

func Parse() error {
	flag.StringVar(&Config.Address, "a", "127.0.0.1:8080", "address to send metrics to")
	flag.DurationVar(&Config.ReportInterval, "r", 10*time.Second, "report interval")
	flag.DurationVar(&Config.PollInterval, "p", 2*time.Second, "poll interval")
	flag.StringVar(&Config.Key, "k", "", "hash key")
	flag.BoolVar(&Config.Profile, "prof", false, "turn on profiling")
	flag.StringVar(&Config.PublicKeyFile, "crypto-key", "", "public key for message encryption")
	flag.Parse()

	err := env.Parse(&Config)
	if err != nil {
		return fmt.Errorf("could not parse config: %w", err)
	}

	log.Info().Msgf("Agent config parsed: %+v", Config)
	return nil
}
