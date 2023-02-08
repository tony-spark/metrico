// Package config contains agent application configuration support (via program arguments and environment variables)
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog/log"
)

var (
	Config = config{
		Address:        "127.0.0.1:8080",
		ReportInterval: 10 * time.Second,
		PollInterval:   2 * time.Second,
	}
)

type config struct {
	Address        string        `short:"a" env:"ADDRESS" json:"address,omitempty" description:"address to send metrics to" `
	ReportInterval time.Duration `short:"r" env:"REPORT_INTERVAL" json:"report_interval,omitempty" description:"report interval" `
	PollInterval   time.Duration `short:"p" env:"POLL_INTERVAL" json:"poll_interval,omitempty" description:"poll interval" `
	Key            string        `short:"k" env:"KEY" json:"key,omitempty" description:"hash key"`
	Profile        bool          `long:"prof" env:"PROFILING" json:"profile,omitempty"  description:"turn on profiling"`
	PublicKeyFile  string        `long:"crypto-key" env:"CRYPTO_KEY" json:"crypto_key,omitempty" description:"public key for message encryption (PEM)"`
	ConfigFile     string        `short:"c" long:"config" env:"CONFIG" description:"configuration file (json)"`
}

type configConfig struct {
	ConfigFile string `short:"c" long:"config" env:"CONFIG" description:"configuration file (json)"`
}

func Parse() error {
	// following mess happend case of config source priorities: config file < cmd args < env
	var confconf configConfig
	_, err := flags.ParseArgs(&confconf, os.Args)
	if err != nil {
		return fmt.Errorf("could not parse c/config option: %w", err)
	}
	err = env.Parse(&confconf)
	if err != nil {
		return fmt.Errorf("could not parse CONFIG env var: %w", err)
	}

	if len(confconf.ConfigFile) > 0 {
		var bs []byte
		bs, err = os.ReadFile(confconf.ConfigFile)
		if err != nil {
			return fmt.Errorf("could not read config file: %w", err)
		}

		err = json.Unmarshal(bs, &Config)
		if err != nil {
			return fmt.Errorf("could not parse JSON config: %w", err)
		}
	}

	_, err = flags.Parse(&Config)
	if err != nil {
		return fmt.Errorf("could not parse flags: %w", err)
	}

	err = env.Parse(&Config)
	if err != nil {
		return fmt.Errorf("could not parse config: %w", err)
	}

	log.Info().Msgf("Agent config parsed: %+v", Config)
	return nil
}

func (c *config) UnmarshalJSON(b []byte) error {
	type configAlias config

	aliasValue := &struct {
		*configAlias
		PollInterval   string `json:"poll_interval,omitempty"`
		ReportInterval string `json:"report_interval,omitempty"`
	}{
		configAlias: (*configAlias)(c),
	}

	err := json.Unmarshal(b, aliasValue)
	if err != nil {
		return fmt.Errorf("could not unmarshal json: %w", err)
	}

	if len(aliasValue.PollInterval) > 0 {
		c.PollInterval, err = time.ParseDuration(aliasValue.PollInterval)
		if err != nil {
			return fmt.Errorf("could not parse time.Duration: %w", err)
		}
	}

	if len(aliasValue.ReportInterval) > 0 {
		c.ReportInterval, err = time.ParseDuration(aliasValue.ReportInterval)
		if err != nil {
			return fmt.Errorf("could not parse time.Duration: %w", err)
		}
	}

	return nil
}
