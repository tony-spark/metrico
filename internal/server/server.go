// Package server contains implementation of metrics server - application to receive and store metrics
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/hash"
	router "github.com/tony-spark/metrico/internal/server/http"
	"github.com/tony-spark/metrico/internal/server/models"
	"github.com/tony-spark/metrico/internal/server/services"
	"github.com/tony-spark/metrico/internal/server/storage"
	"github.com/tony-spark/metrico/internal/server/web"
)

// Server represents server application
type Server struct {
	listenAddress string
	key           string
	dsn           string
	storeFilename string
	storeInterval time.Duration
	restore       bool
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
func WithHTTPServer(listenAddress string) Option {
	return func(s *Server) {
		s.listenAddress = listenAddress
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
// Note that Run blocks until given context is done or error occurred
func (s Server) Run(ctx context.Context) error {
	var r models.MetricRepository
	var dbm models.DBManager
	var postUpdateFn func() = nil
	var err error
	if len(s.dsn) > 0 {
		dbm, err = storage.NewPgManager(s.dsn)
		if err != nil {
			return err
		}
		r = dbm.MetricRepository()
		defer func() {
			err = dbm.Close()
			if err != nil {
				log.Error().Err(err).Msg("error closing database manager")
			}
		}()
	} else {
		var store models.RepositoryPersistence
		r = storage.NewSingleValueRepository()
		store, err = storage.NewJSONFilePersistence(s.storeFilename)
		if err != nil {
			return err
		}
		defer func() {
			err = store.Save(ctx, r)
			if err != nil {
				log.Error().Err(err).Msg("error saving metrics store")
			}
			errc := store.Close()
			if errc != nil {
				log.Error().Err(err).Msg("error closing store")
			}
		}()
		if s.restore {
			err = store.Load(ctx, r)
			if err != nil {
				return err
			}
		}
		pservice := services.NewPersistenceService(store, s.storeInterval, r)
		pservice.Run(ctx)
		postUpdateFn = pservice.PostUpdate()
	}

	var h dto.Hasher
	if len(s.key) > 0 {
		h = hash.NewSha256Hmac(s.key)
	}

	templates := web.NewEmbeddedTemplates()

	return http.ListenAndServe(s.listenAddress,
		router.NewRouter(r, postUpdateFn, h, dbm, templates).R)
}
