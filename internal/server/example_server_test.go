package server_test

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/server"
)

// This example runs server with in-memory storage and shows main endpoints
func TestExample(t *testing.T) {
	tempf, err := os.CreateTemp(os.TempDir(), "metrico-server-example")
	if err != nil {
		defer tempf.Close()
	}
	s := server.New(
		server.WithFileStore(tempf.Name(), 3*time.Second, false),
	)
	go func(s server.Server) {
		err := s.Run(context.Background())
		if err != nil {
			log.Fatal().Err(err).Msg("error while running server")
		}
	}(s)

	go func() {
		// example of sending metric with JSON
		{
			json := `
				{
					"id" : "GaugeExample",
					"type" : "gauge",
					"value" : 1.23
				}
			`
			req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/update/", bytes.NewReader([]byte(json)))
			if err != nil {
				log.Fatal().Err(err).Msg("error while configuring request")
			}
			req.Header.Set("Content-Type", "application/json")
			_, err = http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal().Err(err).Msg("error during request")
			}
		}
	}()

	time.Sleep(30 * time.Second)
}
