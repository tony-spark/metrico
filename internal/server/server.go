// Package server contains implementation of metrics server - application to receive and store metrics
package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/crypto"
	"github.com/tony-spark/metrico/internal/hash"
	router "github.com/tony-spark/metrico/internal/server/http"
	"github.com/tony-spark/metrico/internal/server/models"
	"github.com/tony-spark/metrico/internal/server/services"
	"github.com/tony-spark/metrico/internal/server/storage"
)

// Server represents server application
// TODO: it's just a copy of config so far, rework this to use services
type Server struct {
	listenAddress string
	key           string
	dsn           string
	storeFilename string
	storeInterval time.Duration
	restore       bool
	cryptoKeyFile string
	trustedSubnet string
	dbm           models.DBManager
	store         models.RepositoryPersistence
	r             models.MetricRepository
	srv           *http.Server
}

// Option represents option function for server configuration
type Option func(s *Server)

// New creates server with given options
func New(options ...Option) Server {
	s := Server{
		listenAddress: "127.0.0.1:8080",
		storeFilename: "/tmp/devops-metrics-db.json",
		storeInterval: 300 * time.Second,
		restore:       true,
	}

	for _, opt := range options {
		opt(&s)
	}

	return s
}

// WithHTTPServer configures server to receive metrics via HTTP
func WithHTTPServer(listenAddress string, trustedSubnet string) Option {
	return func(s *Server) {
		s.listenAddress = listenAddress
		s.trustedSubnet = trustedSubnet
	}
}

// WithHashKey configures server to check hash of received messages
func WithHashKey(key string) Option {
	return func(s *Server) {
		s.key = key
	}
}

// WithDB configures server to use database as a metrics storage
func WithDB(dsn string) Option {
	return func(s *Server) {
		s.dsn = dsn
	}
}

func WithCryptoKey(keyFile string) Option {
	return func(s *Server) {
		s.cryptoKeyFile = keyFile
	}
}

// WithFileStore configures server to store metrics in file
//
// # If storeInterval is not specified (0), metrics will be saved on each update
//
// If restore is true, metrics will be loaded
func WithFileStore(filename string, storeInterval time.Duration, restore bool) Option {
	return func(s *Server) {
		s.storeFilename = filename
		s.storeInterval = storeInterval
		s.restore = restore
	}
}

// Run starts a server
//
// Note that Run blocks until Shutdown called
func (s *Server) Run(ctx context.Context) error {
	var postUpdateFn func() = nil
	var err error
	var opts []router.Option
	if len(s.dsn) > 0 {
		s.dbm, err = storage.NewPgManager(s.dsn)
		opts = append(opts, router.WithDBManager(s.dbm))
		if err != nil {
			return err
		}
		s.r = s.dbm.MetricRepository()
	} else {
		s.r = storage.NewSingleValueRepository()
		s.store, err = storage.NewJSONFilePersistence(s.storeFilename)
		if err != nil {
			return err
		}
		if s.restore {
			err = s.store.Load(ctx, s.r)
			if err != nil {
				return err
			}
		}
		pservice := services.NewPersistenceService(s.store, s.storeInterval, s.r)
		pservice.Run(ctx)
		postUpdateFn = pservice.PostUpdate()
	}

	if len(s.key) > 0 {
		h := hash.NewSha256Hmac(s.key)
		opts = append(opts, router.WithHasher(h))
	}

	if len(s.cryptoKeyFile) > 0 {
		var d crypto.Decryptor
		d, err = crypto.NewRSADecryptorFromFile(s.cryptoKeyFile, "metrico")
		if err != nil {
			return fmt.Errorf("could not initialize decryptor: %w", err)
		}
		opts = append(opts, router.WithDecryptor(d))
	}

	if len(s.trustedSubnet) > 0 {
		var subnet *net.IPNet
		_, subnet, err = net.ParseCIDR(s.trustedSubnet)
		if err != nil {
			return fmt.Errorf("could not parse subnet: %w", err)
		}
		opts = append(opts, router.WithTrustedSubNet(subnet))
	}

	metricService := services.NewMetricService(s.r, postUpdateFn)

	rtr := router.NewRouter(metricService, opts...)

	s.srv = &http.Server{
		Addr:    s.listenAddress,
		Handler: rtr.R,
	}

	err = s.srv.ListenAndServe()
	if err != http.ErrServerClosed && err != net.ErrClosed {
		return fmt.Errorf("error running http server: %w", err)
	}
	return nil
}

func (s Server) Shutdown(ctx context.Context) error {
	result := s.srv.Shutdown(ctx)
	if result == nil {
		log.Info().Msg("HTTP server shut down")
	}
	if s.store != nil {
		err := s.store.Save(ctx, s.r)
		if err != nil {
			result = multierror.Append(result, err)
			log.Error().Err(err).Msg("error saving metrics store")
		}
		log.Info().Msg("saved to store")
		errc := s.store.Close()
		if errc != nil {
			result = multierror.Append(result, errc)
			log.Error().Err(err).Msg("error closing store")
		}
		log.Info().Msg("store closed")
	}
	if s.dbm != nil {
		err := s.dbm.Close()
		if err != nil {
			result = multierror.Append(err)
			log.Error().Err(err).Msg("error closing database manager")
		} else {
			log.Info().Msg("database manager closed")
		}
	}
	return result
}
