package server_test

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/server"
	"github.com/tony-spark/metrico/internal/server/services"
	"github.com/tony-spark/metrico/internal/server/storage"
)

// This example runs server with in-memory storage and file persistence (flush every 3 seconds)
func Example() {
	tempf, err := os.CreateTemp(os.TempDir(), "metrico-server-example")
	if err != nil {
		log.Fatal().Err(err).Msg("could not create temp file")
	}
	defer func() {
		if err = tempf.Close(); err != nil {
			log.Fatal().Err(tempf.Close()).Msg("error closing temp file")
		}
	}()
	r := storage.NewSingleValueRepository()
	p, err := storage.NewJSONFilePersistence(tempf.Name())
	pservice := services.NewPersistenceService(p, 3*time.Second, false, r)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create persistence")
	}
	s, err := server.New(
		services.NewMetricService(r, pservice.PostUpdate()),
		server.WithMetricRepository(r),
		server.WithPersistence(pservice),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("could not configure server")
	}
	err = s.Run(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("error while running server")
	}
}
