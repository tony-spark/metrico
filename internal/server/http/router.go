package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/model"
	"github.com/tony-spark/metrico/internal/server/models"
	"github.com/tony-spark/metrico/internal/server/services"
	"github.com/tony-spark/metrico/internal/server/web"
)

type Router struct {
	dbm       models.DBManager
	ms        *services.MetricService
	R         chi.Router
	h         dto.Hasher
	templates web.TemplateProvider
}

func NewRouter(repo models.MetricRepository, postUpdateFn func(), h dto.Hasher, dbm models.DBManager, templates web.TemplateProvider) *Router {
	r := chi.NewRouter()

	router := &Router{
		dbm:       dbm,
		ms:        services.NewMetricService(repo, postUpdateFn),
		R:         r,
		h:         h,
		templates: templates,
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httplog.RequestLogger(log.Logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Mount("/debug", middleware.Profiler())

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
			r.Get("/{name}", router.MetricGetHandler(model.COUNTER))
		})
		r.Route("/gauge", func(r chi.Router) {
			r.Get("/{name}", router.MetricGetHandler(model.GAUGE))
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
