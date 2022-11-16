package server

import (
	"context"
	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/hash"
	router "github.com/tony-spark/metrico/internal/server/http"
	"github.com/tony-spark/metrico/internal/server/models"
	"github.com/tony-spark/metrico/internal/server/services"
	"github.com/tony-spark/metrico/internal/server/storage"
	"net/http"
	"time"
)

type Server struct {
	listenAddress string
	key           string
	dsn           string
	storeFilename string
	storeInterval time.Duration
	restore       bool
}

type Option func(s *Server)

func New(options ...Option) Server {
	s := Server{
		listenAddress: "127.0.0.1:8080",
		storeFilename: "/tmp/devops-metrics-db.json",
		storeInterval: 300 * time.Second,
		restore:       true,
	}

	for _, opt := range options {
		opt(&s)
	}

	return s
}

func WithHTTPServer(listenAddress string) Option {
	return func(s *Server) {
		s.listenAddress = listenAddress
	}
}

func WithHashKey(key string) Option {
	return func(s *Server) {
		s.key = key
	}
}

func WithDB(dsn string) Option {
	return func(s *Server) {
		s.dsn = dsn
	}
}

func WithFileStore(filename string, storeInterval time.Duration, restore bool) Option {
	return func(s *Server) {
		s.storeFilename = filename
		s.storeInterval = storeInterval
		s.restore = restore
	}
}

// Run starts a server
func (s Server) Run(ctx context.Context) error {
	var gr models.GaugeRepository
	var cr models.CounterRepository
	var dbm models.DBManager
	var postUpdateFn func() = nil
	var err error
	if len(s.dsn) > 0 {
		dbm, err = storage.NewPgManager(s.dsn)
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
		store, err = storage.NewJSONFilePersistence(s.storeFilename)
		if err != nil {
			return err
		}
		defer func() {
			store.Save(ctx, gr, cr)
			store.Close()
		}()
		if s.restore {
			err = store.Load(ctx, gr, cr)
			if err != nil {
				return err
			}
		}
		pservice := services.NewPersistenceService(store, s.storeInterval, gr, cr)
		pservice.Run(ctx)
		postUpdateFn = pservice.PostUpdate()
	}

	var h dto.Hasher
	if len(s.key) > 0 {
		h = hash.NewSha256Hmac(s.key)
	}

	return http.ListenAndServe(s.listenAddress,
		router.NewRouter(gr, cr, postUpdateFn, h, dbm).R)
}
