package server

import (
	"github.com/tony-spark/metrico/internal/server/handlers"
	"github.com/tony-spark/metrico/internal/server/models"
	"github.com/tony-spark/metrico/internal/server/storage"
	"net/http"
)

// Run starts a server for collecting metrics using HTTP API
//
// HTTP server listens bindAddress
func Run(bindAddress string) error {
	var gaugeRepo models.GaugeRepository = storage.NewSingleValueGaugeRepository()
	var counterRepo models.CounterRepository = storage.NewSingleValueCounterRepository()

	http.HandleFunc("/", handlers.DefaultHandler)
	http.HandleFunc("/update/counter/", handlers.CounterHandler(counterRepo))
	http.HandleFunc("/update/gauge/", handlers.GaugeHandler(gaugeRepo))

	return http.ListenAndServe(bindAddress, nil)
}
