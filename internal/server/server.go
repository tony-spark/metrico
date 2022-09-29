package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tony-spark/metrico/internal/server/handlers"
	"github.com/tony-spark/metrico/internal/server/models"
	"github.com/tony-spark/metrico/internal/server/storage"
	"net/http"
)

func NewRouter(gaugeRepo models.GaugeRepository, counterRepo models.CounterRepository) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/update/", func(r chi.Router) {
		r.Route("/counter", func(r chi.Router) {
			r.Get("/{name}", handlers.CounterGetHandler(counterRepo))
			r.Post("/{name}/{svalue}", handlers.CounterPostHandler(counterRepo))
		})
		r.Route("/gauge", func(r chi.Router) {
			r.Get("/{name}", handlers.GaugeGetHandler(gaugeRepo))
			r.Post("/{name}/{svalue}", handlers.GaugePostHandler(gaugeRepo))
		})
		r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
			mtype := chi.URLParam(r, "*")
			http.Error(w, "unknown metric type in "+mtype, http.StatusNotImplemented)
		})
	})

	return r
}

// Run starts a server for collecting metrics using HTTP API
//
// HTTP server listens bindAddress
func Run(bindAddress string) error {
	return http.ListenAndServe(bindAddress,
		NewRouter(storage.NewSingleValueGaugeRepository(), storage.NewSingleValueCounterRepository()))
}
