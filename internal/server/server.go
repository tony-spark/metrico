// Package server contains implementation of metrics server - application to receive and store metrics
package server

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/server/models"
	"github.com/tony-spark/metrico/internal/server/services"
	"golang.org/x/sync/errgroup"
)

type Controller interface {
	fmt.Stringer

	Run() error
	Shutdown(ctx context.Context) error
}

// Server represents server application
type Server struct {
	dbm      models.DBManager
	store    models.RepositoryPersistence
	r        models.MetricRepository
	pService *services.PersistenceService
	mService *services.MetricService
	ctrls    []Controller
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

func AddController(c Controller) Option {
	return func(s *Server) {
		s.ctrls = append(s.ctrls, c)
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

	grp := new(errgroup.Group)
	for _, ctrl := range s.ctrls {
		ctrl := ctrl
		grp.Go(ctrl.Run)
	}

	return grp.Wait()
}

func (s Server) Shutdown(ctx context.Context) error {
	var result error
	for _, ctrl := range s.ctrls {
		log.Info().Msgf("shutting down %v", ctrl)
		err := ctrl.Shutdown(ctx)
		if err != nil {
			result = multierror.Append(result, err)
		}
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
