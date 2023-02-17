// Package server contains implementation of metrics server - application to receive and store metrics
package server

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"
	router "github.com/tony-spark/metrico/internal/server/http"
	"github.com/tony-spark/metrico/internal/server/models"
	"github.com/tony-spark/metrico/internal/server/services"
)

// Server represents server application
// TODO: it's just a copy of config so far, rework this to use services
type Server struct {
	dbm            models.DBManager
	store          models.RepositoryPersistence
	r              models.MetricRepository
	pService       *services.PersistenceService
	mService       *services.MetricService
	httpController *router.Controller
}

// Option represents option function for server configuration
type Option func(s *Server)

func WithDBManager(dbm models.DBManager) Option {
	return func(s *Server) {
		s.dbm = dbm
	}
}

func WithMetricRepository(r models.MetricRepository) Option {
	return func(s *Server) {
		s.r = r
	}
}

func WithPersistence(pservice *services.PersistenceService) Option {
	return func(s *Server) {
		s.pService = pservice
	}
}

func WithHTTPController(c *router.Controller) Option {
	return func(s *Server) {
		s.httpController = c
	}
}

// New creates server with given options
func New(ms *services.MetricService, options ...Option) (Server, error) {
	s := Server{
		mService: ms,
	}

	for _, opt := range options {
		opt(&s)
	}

	return s, nil
}

// Run starts a server
//
// Note that Run blocks until Shutdown called
func (s *Server) Run(ctx context.Context) error {
	if s.pService != nil {
		err := s.pService.Run(ctx)
		if err != nil {
			return err
		}
	}

	return s.httpController.Run()
}

func (s Server) Shutdown(ctx context.Context) error {
	result := s.httpController.Shutdown(ctx)
	if result == nil {
		log.Info().Msg("HTTP controller shut down")
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
