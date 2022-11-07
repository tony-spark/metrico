package services

import (
	"context"
	"github.com/tony-spark/metrico/internal/server/models"
	"time"
)

type PersistenceService struct {
	p             models.RepositoryPersistence
	storeInterval time.Duration
	gr            models.GaugeRepository
	cr            models.CounterRepository
	postUpdate    func()
}

func NewPersistenceService(p models.RepositoryPersistence, storeInterval time.Duration,
	gr models.GaugeRepository, cr models.CounterRepository) *PersistenceService {
	return &PersistenceService{
		p:             p,
		storeInterval: storeInterval,
		gr:            gr,
		cr:            cr,
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
					s.p.Save(ctx, s.gr, s.cr)
				case <-ctx.Done():
					return
				}
			}
		}()
	} else {
		s.postUpdate = func() {
			s.p.Save(ctx, s.gr, s.cr)
		}
	}
}

func (s PersistenceService) PostUpdate() func() {
	return s.postUpdate
}
