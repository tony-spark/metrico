// Package config contains agent application configuration support (via program arguments and environment variables)
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
	configUtil "github.com/tony-spark/metrico/internal/config"
)

var (
	Config = config{
		Address:        "127.0.0.1:8080",
		ReportInterval: 10 * time.Second,
		PollInterval:   2 * time.Second,
	}
)

type config struct {
	Address        string        `env:"ADDRESS" json:"address,omitempty"`
	GrpcAddress    string        `env:"GRPC_ADRESS" json:"grpc_address,omitempty"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" json:"report_interval,omitempty"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" json:"poll_interval,omitempty"`
	Key            string        `env:"KEY" json:"key,omitempty"`
	Profile        bool          `env:"PROFILING" json:"profile,omitempty"`
	PublicKeyFile  string        `env:"CRYPTO_KEY" json:"crypto_key,omitempty"`
}

func Parse() error {
	// the following mess happens cause of config source priorities: config file < cmd args < env
	configFile, err := configUtil.ParseConfigFileParameter(os.Args)
	if err != nil {
		return fmt.Errorf("coult not parse config file parameter: %w", err)
	}

	if len(configFile) > 0 {
		var bs []byte
		bs, err = os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("could not read config file: %w", err)
		}

		err = json.Unmarshal(bs, &Config)
		if err != nil {
			return fmt.Errorf("could not parse JSON config: %w", err)
		}
	}

	flag.StringVar(&Config.Address, "a", Config.Address, "address to send metrics to")
	flag.StringVar(&Config.GrpcAddress, "g", Config.GrpcAddress, "user grpc with specified address instead")
	flag.DurationVar(&Config.ReportInterval, "r", Config.ReportInterval, "report interval")
	flag.DurationVar(&Config.PollInterval, "p", Config.PollInterval, "poll interval")
	flag.StringVar(&Config.Key, "k", Config.Key, "hash key")
	flag.BoolVar(&Config.Profile, "prof", Config.Profile, "turn on profiling")
	flag.StringVar(&configFile, "config", "", "config file")
	flag.StringVar(&configFile, "c", "", "shortcut to --config")
	flag.Parse()

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
