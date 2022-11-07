package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/server/models"
)

type Router struct {
	R          chi.Router
	gr         models.GaugeRepository
	cr         models.CounterRepository
	postUpdate func()
	h          dto.Hasher
	dbr        models.DBManager
}

func NewRouter(gaugeRepo models.GaugeRepository, counterRepo models.CounterRepository, postUpdateFn func(), h dto.Hasher, dbm models.DBManager) *Router {
	r := chi.NewRouter()

	router := &Router{
		R:          r,
		gr:         gaugeRepo,
		cr:         counterRepo,
		postUpdate: postUpdateFn,
		h:          h,
		dbr:        dbm,
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httplog.RequestLogger(log.Logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Get("/", router.MetricsViewPageHandler())
	r.Route("/update", func(r chi.Router) {
		r.Route("/counter", func(r chi.Router) {
			r.Post("/{name}/{svalue}", router.CounterPostHandler())
		})
		r.Route("/gauge", func(r chi.Router) {
			r.Post("/{name}/{svalue}", router.GaugePostHandler())
		})
		r.Post("/", router.UpdatePostHandler())
		r.HandleFunc("/*", handleUnknown)
	})
	r.Route("/value", func(r chi.Router) {
		r.Route("/counter", func(r chi.Router) {
			r.Get("/{name}", router.CounterGetHandler())
		})
		r.Route("/gauge", func(r chi.Router) {
			r.Get("/{name}", router.GaugeGetHandler())
		})
		r.Post("/", router.GetPostHandler())
		r.HandleFunc("/*", handleUnknown)
	})
	r.Get("/ping", router.PingHandler())
	r.Route("/updates", func(r chi.Router) {
		r.Post("/", router.BulkUpdatePostHandler())
	})

	return router
}
