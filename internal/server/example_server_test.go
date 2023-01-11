package server_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/server"
)

// This example runs server with in-memory storage and shows main endpoints
func Example() {
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

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
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
		// example of sending metrics in batch
		{
			json := `
				[{
					"id" : "GaugeExample",
					"type" : "gauge",
					"value" : 1.5
				},
				{
					"id" : "CounterExample",
					"type" : "counter",
					"delta" : 4
				}]
			`
			req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/updates/", bytes.NewReader([]byte(json)))
			if err != nil {
				log.Fatal().Err(err).Msg("error while configuring request")
			}
			req.Header.Set("Content-Type", "application/json")
			_, err = http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal().Err(err).Msg("error during request")
			}
		}
		// example of getting metrics value
		{
			json := `
				{
					"id" : "GaugeExample",
					"type" : "gauge"
				}
			`
			req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/value/", bytes.NewReader([]byte(json)))
			if err != nil {
				log.Fatal().Err(err).Msg("error while configuring request")
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal().Err(err).Msg("error during request")
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal().Err(err).Msg("could not read response body")
			}
			log.Info().Msg(string(body))
		}
	}()

	wg.Wait()
}
