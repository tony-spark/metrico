package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/tony-spark/metrico/internal/server/models"
	"log"
	"net/http"
	"strconv"
)

func CounterGetHandler(repo models.CounterRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		counter, err := repo.GetByName(name)
		if err != nil {
			http.Error(w, "error retrieving value", http.StatusInternalServerError)
			return
		}
		if counter == nil {
			http.Error(w, "counter not found", http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprint(counter.Value)))
	}
}

func CounterPostHandler(repo models.CounterRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		svalue := chi.URLParam(r, "svalue")

		value, err := strconv.ParseInt(svalue, 10, 64)
		if err != nil {
			http.Error(w, "VALUE type must be int64", http.StatusBadRequest)
			return
		}
		_, err = repo.AddAndSave(name, value)
		if err != nil {
			log.Println("Could not add and save counter value", name, value)
			http.Error(w, "Could not add and save counter value", http.StatusInternalServerError)
			return
		}
	}
}

func GaugeGetHandler(repo models.GaugeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		gauge, err := repo.GetByName(name)
		if err != nil {
			http.Error(w, "error retrieving value", http.StatusInternalServerError)
			return
		}
		if gauge == nil {
			http.Error(w, "value not found", http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprint(gauge.Value)))
	}
}

func GaugePostHandler(repo models.GaugeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		svalue := chi.URLParam(r, "svalue")

		value, err := strconv.ParseFloat(svalue, 64)
		if err != nil {
			http.Error(w, "VALUE type must be float64", http.StatusBadRequest)
			return
		}
		_, err = repo.Save(name, value)
		if err != nil {
			log.Println("Could not save gauge value", name, value)
			http.Error(w, "Could not save gauge value", http.StatusInternalServerError)
			return
		}
	}
}
