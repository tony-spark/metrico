package server

import (
	"context"
	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/hash"
	"github.com/tony-spark/metrico/internal/server/config"
	router "github.com/tony-spark/metrico/internal/server/http"
	"github.com/tony-spark/metrico/internal/server/models"
	"github.com/tony-spark/metrico/internal/server/services"
	"github.com/tony-spark/metrico/internal/server/storage"
	"net/http"
)

// Run starts a server for collecting metrics using HTTP API
//
// HTTP server listens bindAddress
func Run(ctx context.Context) error {
	var gr models.GaugeRepository
	var cr models.CounterRepository
	var dbm models.DBManager
	var postUpdateFn func() = nil
	var err error
	if len(config.Config.DSN) > 0 {
		dbm, err = storage.NewPgManager(config.Config.DSN)
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
		store, err = storage.NewJSONFilePersistence(config.Config.StoreFilename)
		if err != nil {
			return err
		}
		defer func() {
			store.Save(ctx, gr, cr)
			store.Close()
		}()
		if config.Config.Restore {
			err = store.Load(ctx, gr, cr)
			if err != nil {
				return err
			}
		}
		pservice := services.NewPersistenceService(store, config.Config.StoreInterval, gr, cr)
		pservice.Run(ctx)
		postUpdateFn = pservice.PostUpdate()
	}

	var h dto.Hasher
	if len(config.Config.Key) > 0 {
		h = hash.NewSha256Hmac(config.Config.Key)
	}

	return http.ListenAndServe(config.Config.Address,
		router.NewRouter(gr, cr, postUpdateFn, h, dbm).R)
}
