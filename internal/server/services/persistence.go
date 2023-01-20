package services

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/tony-spark/metrico/internal/server/models"
)

type PersistenceService struct {
	p             models.RepositoryPersistence
	storeInterval time.Duration
	r             models.MetricRepository
	postUpdate    func()
}

func NewPersistenceService(p models.RepositoryPersistence, storeInterval time.Duration,
	r models.MetricRepository) *PersistenceService {
	return &PersistenceService{
		p:             p,
		storeInterval: storeInterval,
		r:             r,
	}
}

func (s *PersistenceService) Run(ctx context.Context) {
	if s.storeInterval > 0 {
		saveTicker := time.NewTicker(s.storeInterval)
		defer saveTicker.Stop()
		go func() {
			for {
				select {
				case <-saveTicker.C:
					err := s.p.Save(ctx, s.r)
					if err != nil {
						log.Error().Err(err).Msg("Could not persist metrics")
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	} else {
		s.postUpdate = func() {
			err := s.p.Save(ctx, s.r)
			if err != nil {
				log.Error().Err(err).Msg("could not persist metrics")
			}
		}
	}
}

func (s PersistenceService) PostUpdate() func() {
	return s.postUpdate
}
