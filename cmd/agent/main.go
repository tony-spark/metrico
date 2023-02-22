// Package main contains entrypoint for agent application
package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/agent/transports"
	grpcTransport "github.com/tony-spark/metrico/internal/agent/transports/grpc"
	httpTransport "github.com/tony-spark/metrico/internal/agent/transports/http"
	"github.com/tony-spark/metrico/internal/crypto"
	"github.com/tony-spark/metrico/internal/hash"

	a "github.com/tony-spark/metrico/internal/agent"
	"github.com/tony-spark/metrico/internal/agent/config"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	err := config.Parse()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not parse config")
	}

	var t transports.Transport

	if len(config.Config.Address) == 0 && len(config.Config.GrpcAddress) == 0 {
		log.Fatal().Msg("you should define address (grpc or http)")
	}

	if len(config.Config.Address) > 0 {
		baseURL := "http://" + strings.Trim(config.Config.Address, "\"")

		var options []httpTransport.Option
		if len(config.Config.Key) > 0 {
			options = append(options, httpTransport.WithHasher(hash.NewSha256Hmac(config.Config.Key)))
		}
		if len(config.Config.PublicKeyFile) > 0 {
			encryptor, err := crypto.NewRSAEncryptorFromFile(config.Config.PublicKeyFile, "metrico")
			if err != nil {
				log.Fatal().Err(err).Msg("could not parse public key")
			}
			options = append(options, httpTransport.WithEncryptor(encryptor))
		}
		t = httpTransport.NewTransport(baseURL, options...)
	}

	if len(config.Config.GrpcAddress) > 0 {
		var options []grpcTransport.Option
		if len(config.Config.Key) > 0 {
			options = append(options, grpcTransport.WithHasher(hash.NewSha256Hmac(config.Config.Key)))
		}
		t, err = grpcTransport.NewTransport(config.Config.GrpcAddress, options...)
		if err != nil {
			log.Fatal().Err(err).Msg("could not initialize grpc transport")
		}
	}

	agent := a.New(
		a.WithTransport(t),
		a.WithPollInterval(config.Config.PollInterval),
		a.WithReportInterval(config.Config.ReportInterval),
	)

	ctx, cancel := context.WithCancel(context.Background())
	go agent.Run(ctx)

	if config.Config.Profile {
		log.Info().Msg("starting profile http server")
		go func() {
			err := http.ListenAndServe("127.0.0.1:8888", nil)
			if err != nil {
				log.Error().Err(err).Msg("error running http server")
			}
		}()
	}

	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-terminateSignal
	cancel()

	log.Info().Msg("shutting down gracefully...")

	agent.Stop()

	log.Info().Msg("agent shut down")
}
