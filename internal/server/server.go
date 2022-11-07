package server

import (
	"context"
	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/hash"
	"github.com/tony-spark/metrico/internal/server/config"
	router "github.com/tony-spark/metrico/internal/server/http"
	"github.com/tony-spark/metrico/internal/server/models"
	"github.com/tony-spark/metrico/internal/server/storage"
	"net/http"
	"time"
)

// Run starts a server for collecting metrics using HTTP API
//
// HTTP server listens bindAddress
func Run(ctx context.Context, cfg config.Config) error {
	var gr models.GaugeRepository
	var cr models.CounterRepository
	var dbm models.DBManager
	var postUpdateFn func() = nil
	var err error
	if len(cfg.DSN) > 0 {
		dbm, err = storage.NewPgManager(cfg.DSN)
		if err != nil {
			return err
		}
		gr = dbm.GaugeRepository()
		cr = dbm.CounterRepository()
		defer dbm.Close()
	} else {
		var store models.RepositoryPersistence
		gr = storage.NewSingleValueGaugeRepository()
		cr = storage.NewSingleValueCounterRepository()
		store, err = storage.NewJSONFilePersistence(cfg.StoreFilename)
		if err != nil {
			return err
		}
		defer func() {
			store.Save(ctx, gr, cr)
			store.Close()
		}()
		if cfg.Restore {
			err = store.Load(ctx, gr, cr)
			if err != nil {
				return err
			}
		}
		// TODO: simplify code (extract ticker logic to service?)
		if cfg.StoreInterval > 0 {
			saveTicker := time.NewTicker(cfg.StoreInterval)
			defer saveTicker.Stop()
			go func() {
				for {
					select {
					case <-saveTicker.C:
						store.Save(ctx, gr, cr)
					case <-ctx.Done():
						return
					}
				}
			}()
		} else {
			postUpdateFn = func() {
				store.Save(ctx, gr, cr)
			}
		}
	}

	var h dto.Hasher
	if len(cfg.Key) > 0 {
		h = hash.NewSha256Hmac(cfg.Key)
	}

	return http.ListenAndServe(cfg.Address,
		router.NewRouter(gr, cr, postUpdateFn, h, dbm).R)
}
