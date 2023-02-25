// Package main contains entrypoint for server application
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/crypto"
	"github.com/tony-spark/metrico/internal/hash"
	grpcController "github.com/tony-spark/metrico/internal/server/grpc"
	httpController "github.com/tony-spark/metrico/internal/server/http"
	"github.com/tony-spark/metrico/internal/server/models"
	"github.com/tony-spark/metrico/internal/server/services"
	"github.com/tony-spark/metrico/internal/server/storage"

	"github.com/tony-spark/metrico/internal/server"
	"github.com/tony-spark/metrico/internal/server/config"
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

	ctx, cancel := context.WithCancel(context.Background())

	var postUpdateFn func() = nil
	var r models.MetricRepository
	httpCtrlOpts := []httpController.Option{
		httpController.WithListenAddress(config.Config.Address),
	}
	grpcCtrlOpts := []grpcController.Option{
		grpcController.WithListenAddress(config.Config.GrpcAddress),
	}
	var serverOpts []server.Option
	if len(config.Config.DSN) > 0 {
		var dbm models.DBManager
		dbm, err = storage.NewPgManager(config.Config.DSN)
		if err != nil {
			log.Fatal().Err(err).Msg("could not create DB manager")
		}
		r = dbm.MetricRepository()
		serverOpts = append(serverOpts, server.WithDBManager(dbm), server.WithMetricRepository(r))
		httpCtrlOpts = append(httpCtrlOpts, httpController.WithDBManager(dbm))
		grpcCtrlOpts = append(grpcCtrlOpts, grpcController.WithDBManager(dbm))
	} else {
		var p models.RepositoryPersistence
		r = storage.NewSingleValueRepository()
		serverOpts = append(serverOpts, server.WithMetricRepository(r))
		p, err = storage.NewJSONFilePersistence(config.Config.StoreFilename)
		if err != nil {
			log.Fatal().Err(err).Msg("could not create persistence")
		}
		pservice := services.NewPersistenceService(p, config.Config.StoreInterval, config.Config.Restore, r)
		serverOpts = append(serverOpts, server.WithPersistence(pservice))
		postUpdateFn = pservice.PostUpdate()
	}

	if len(config.Config.Key) > 0 {
		h := hash.NewSha256Hmac(config.Config.Key)
		httpCtrlOpts = append(httpCtrlOpts, httpController.WithHasher(h))
		grpcCtrlOpts = append(grpcCtrlOpts, grpcController.WithHasher(h))
	}

	if len(config.Config.PrivateKeyFile) > 0 {
		var d crypto.Decryptor
		d, err = crypto.NewRSADecryptorFromFile(config.Config.PrivateKeyFile, "metrico")
		if err != nil {
			log.Fatal().Err(err).Msg("could not initialize decryptor")
		}
		httpCtrlOpts = append(httpCtrlOpts, httpController.WithDecryptor(d))
	}

	if len(config.Config.TrustedSubnet) > 0 {
		var subnet *net.IPNet
		_, subnet, err = net.ParseCIDR(config.Config.TrustedSubnet)
		if err != nil {
			log.Fatal().Err(err).Msg("could not parse subnet")
		}
		httpCtrlOpts = append(httpCtrlOpts, httpController.WithTrustedSubNet(subnet))
		grpcCtrlOpts = append(grpcCtrlOpts, grpcController.WithTrustedSubNet(subnet))
	}

	metricService := services.NewMetricService(r, postUpdateFn)

	serverOpts = append(serverOpts, server.AddController(httpController.NewController(metricService, httpCtrlOpts...)))

	if len(config.Config.GrpcAddress) > 0 {
		serverOpts = append(serverOpts, server.AddController(grpcController.NewController(metricService, grpcCtrlOpts...)))
	}

	s, err := server.New(metricService, serverOpts...)

	if err != nil {
		log.Fatal().Err(err).Msg("could not configure server")
	}

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	shutdownDone := make(chan struct{})

	go func() {
		<-shutdownSignal

		log.Info().Msg("shutting down gracefully...")
		err = s.Shutdown(context.Background())
		if err != nil {
			log.Fatal().Err(err).Msg("failed to shut down gracefully")
		}

		cancel()
		close(shutdownDone)
	}()

	err = s.Run(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error running server")
	}

	<-shutdownDone
	log.Info().Msg("server shut down gracefully")
}
