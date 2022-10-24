package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tony-spark/metrico/internal/server/handlers"
	"github.com/tony-spark/metrico/internal/server/models"
	"github.com/tony-spark/metrico/internal/server/storage"
	"net/http"
	"time"
)

func NewRouter(gaugeRepo models.GaugeRepository, counterRepo models.CounterRepository, postUpdateFn func()) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", handlers.PageHandler(gaugeRepo, counterRepo))
	r.Route("/update", func(r chi.Router) {
		r.Route("/counter", func(r chi.Router) {
			r.Post("/{name}/{svalue}", handlers.CounterPostHandler(counterRepo, postUpdateFn))
		})
		r.Route("/gauge", func(r chi.Router) {
			r.Post("/{name}/{svalue}", handlers.GaugePostHandler(gaugeRepo, postUpdateFn))
		})
		r.Post("/", handlers.UpdatePostHandler(gaugeRepo, counterRepo, postUpdateFn))
		r.HandleFunc("/*", handleUnknown)
	})
	r.Route("/value", func(r chi.Router) {
		r.Route("/counter", func(r chi.Router) {
			r.Get("/{name}", handlers.CounterGetHandler(counterRepo))
		})
		r.Route("/gauge", func(r chi.Router) {
			r.Get("/{name}", handlers.GaugeGetHandler(gaugeRepo))
		})
		r.Post("/", handlers.GetPostHandler(gaugeRepo, counterRepo))
		r.HandleFunc("/*", handleUnknown)
	})

	return r
}

// Run starts a server for collecting metrics using HTTP API
//
// HTTP server listens bindAddress
func Run(ctx context.Context, bindAddress string, storeFilename string, restore bool, storeInterval time.Duration) error {
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
	return http.ListenAndServe(bindAddress,
		NewRouter(gr, cr, postUpdateFn))
}

func handleUnknown(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "*")
	http.Error(w, "unknown metric type in "+mtype, http.StatusNotImplemented)
}
