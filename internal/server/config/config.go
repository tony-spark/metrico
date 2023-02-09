// Package config contains server application configuration support (via program arguments and environment variables)
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
		Address:       "127.0.0.1:8080",
		StoreInterval: 300 * time.Second,
		StoreFilename: "/tmp/devops-metrics-db.json",
		Restore:       true,
	}
)

type config struct {
	Address        string        `env:"ADDRESS" json:"address,omitempty"`
	StoreInterval  time.Duration `env:"STORE_INTERVAL" json:"store_interval,omitempty"`
	StoreFilename  string        `env:"STORE_FILE" json:"store_filename,omitempty"`
	Restore        bool          `env:"RESTORE" json:"restore,omitempty"`
	Key            string        `env:"KEY" json:"key,omitempty"`
	DSN            string        `env:"DATABASE_DSN" json:"database_dsn,omitempty"`
	PrivateKeyFile string        `env:"CRYPTO_KEY" json:"crypto_key,omitempty"`
	TrustedSubnet  string        `env:"TRUSTED_SUBNET" json:"trusted_subnet,omitempty"`
}

func Parse() error {
	// the following mess happens cause of config source priorities: config file < cmd args < env
	configFile, err := configUtil.ParseConfigFileParameter()
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

	flag.StringVar(&Config.Address, "a", Config.Address, "address to listen")
	flag.DurationVar(&Config.StoreInterval, "i", Config.StoreInterval, "store interval")
	flag.StringVar(&Config.StoreFilename, "f", Config.StoreFilename, "file to persist metrics")
	flag.BoolVar(&Config.Restore, "r", Config.Restore, "whether to load metric from file on start")
	flag.StringVar(&Config.Key, "k", Config.Key, "hash key")
	flag.StringVar(&Config.DSN, "d", Config.DSN, "database connection string")
	flag.StringVar(&Config.PrivateKeyFile, "crypto-key", Config.PrivateKeyFile, "private key for message decryption (PEM)")
	flag.StringVar(&Config.TrustedSubnet, "t", Config.TrustedSubnet, "trusted subnet for clients")
	flag.StringVar(&configFile, "config", "", "config file")
	flag.StringVar(&configFile, "c", "", "shortcut to --config")
	flag.Parse()

	err = env.Parse(&Config)
	if err != nil {
		return fmt.Errorf("could not read config: %w", err)
	}

	log.Info().Msgf("Server config parsed:  %+v", Config)
	return nil
}

func (c *config) UnmarshalJSON(b []byte) error {
	type configAlias config

	aliasValue := &struct {
		*configAlias
		StoreInterval string `json:"store_interval,omitempty"`
	}{
		configAlias: (*configAlias)(c),
	}

	err := json.Unmarshal(b, aliasValue)
	if err != nil {
		return fmt.Errorf("could not unmarshal json: %w", err)
	}

	if len(aliasValue.StoreInterval) > 0 {
		c.StoreInterval, err = time.ParseDuration(aliasValue.StoreInterval)
		if err != nil {
			return fmt.Errorf("could not parse time.Duration: %w", err)
		}
	}

	return nil
}
