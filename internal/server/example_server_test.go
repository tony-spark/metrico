package server_test

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/server"
)

// This example runs server with in-memory storage and file persistence (flush every 3 seconds)
func Example() {
	tempf, err := os.CreateTemp(os.TempDir(), "metrico-server-example")
	if err != nil {
		log.Fatal().Err(err).Msg("could not create temp file")
	}
	defer tempf.Close()
	s := server.New(
		server.WithFileStore(tempf.Name(), 3*time.Second, false),
	)
	err = s.Run(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("error while running server")
	}
}
