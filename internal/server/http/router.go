// Package http contains HTTP API implementation to handle metrics. See swagger specification for more details
package http

import (
	"context"
	"fmt"
	"net"
	"net/http"

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
	listenAddress string
	srv           *http.Server
	R             chi.Router
	ms            *services.MetricService
	templates     web.TemplateProvider
	dbm           models.DBManager
	h             dto.Hasher
	d             crypto.Decryptor
	trustedSubNet *net.IPNet
}

type Option func(r *Router)

func WithListenAddress(addr string) Option {
	return func(r *Router) {
		r.listenAddress = addr
	}
}

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

func NewRouter(metricService *services.MetricService, options ...Option) *Router {
	r := chi.NewRouter()

	router := &Router{
		ms:        metricService,
		R:         r,
		templates: web.NewEmbeddedTemplates(),
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

func (router *Router) Run() error {
	router.srv = &http.Server{
		Addr:    router.listenAddress,
		Handler: router.R,
	}

	err := router.srv.ListenAndServe()
	if err != http.ErrServerClosed && err != net.ErrClosed {
		return fmt.Errorf("error running http server: %w", err)
	}

	return nil
}

func (router Router) Shutdown(ctx context.Context) error {
	return router.srv.Shutdown(ctx)
}
