// Package http contains HTTP API implementation to handle metrics. See swagger specification for more details
package http

import (
	"net"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/crypto"

	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/model"
	"github.com/tony-spark/metrico/internal/server/models"
	"github.com/tony-spark/metrico/internal/server/services"
	"github.com/tony-spark/metrico/internal/server/web"
)

// @Title Metric API
// @Description Metric storage
// @Version 1.0

type Router struct {
	R             chi.Router
	ms            *services.MetricService
	templates     web.TemplateProvider
	dbm           models.DBManager
	h             dto.Hasher
	d             crypto.Decryptor
	trustedSubNet *net.IPNet
}

type Option func(r *Router)

func WithHasher(h dto.Hasher) Option {
	return func(r *Router) {
		r.h = h
	}
}

func WithDBManager(dbm models.DBManager) Option {
	return func(r *Router) {
		r.dbm = dbm
	}
}

func WithDecryptor(d crypto.Decryptor) Option {
	return func(r *Router) {
		r.d = d
	}
}

func WithTrustedSubNet(subnet *net.IPNet) Option {
	return func(r *Router) {
		r.trustedSubNet = subnet
	}
}

func NewRouter(metricService *services.MetricService, templates web.TemplateProvider, options ...Option) *Router {
	r := chi.NewRouter()

	router := &Router{
		ms:        metricService,
		R:         r,
		templates: templates,
	}

	for _, opt := range options {
		opt(router)
	}

	r.Use(middleware.RealIP)
	if router.trustedSubNet != nil {
		r.Use(SubnetClientFilter(*router.trustedSubNet))
	}
	r.Use(middleware.RequestID)
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
