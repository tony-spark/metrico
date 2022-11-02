package server

import (
	"context"
	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/hash"
	router "github.com/tony-spark/metrico/internal/server/http"
	"github.com/tony-spark/metrico/internal/server/storage"
	"net/http"
	"time"
)

// Run starts a server for collecting metrics using HTTP API
//
// HTTP server listens bindAddress
func Run(ctx context.Context, bindAddress string, storeFilename string, restore bool, storeInterval time.Duration, key string) error {
	gr := storage.NewSingleValueGaugeRepository()
	cr := storage.NewSingleValueCounterRepository()
	store, err := storage.NewJSONFilePersistence(storeFilename)
	defer func() {
		store.Save(gr, cr)
		store.Close()
	}()
	if err != nil {
		return err
	}
	if restore {
		err = store.Load(gr, cr)
		if err != nil {
			return err
		}
	}
	var postUpdateFn func() = nil
	// TODO: simplify code (extract ticker logic to service?)
	if storeInterval > 0 {
		saveTicker := time.NewTicker(storeInterval)
		defer saveTicker.Stop()
		go func() {
			for {
				select {
				case <-saveTicker.C:
					store.Save(gr, cr)
				case <-ctx.Done():
					return
				}
			}
		}()
	} else {
		postUpdateFn = func() {
			store.Save(gr, cr)
		}
	}
	var h dto.Hasher
	if len(key) > 0 {
		h = hash.NewSha256Keyed(key)
	}
	return http.ListenAndServe(bindAddress,
		router.NewRouter(gr, cr, postUpdateFn, h).R)
}
