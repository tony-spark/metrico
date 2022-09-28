package server

import (
	"github.com/tony-spark/metrico/internal/server/models"
	"io"
	"log"
	"net/http"
	"strconv"
	s "strings"
)

// Listen update/counter/*, process update/counter/<NAME>/<VALUE> (value int64)
func counterHandler(repo models.CounterRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ss := s.Split(r.URL.Path, "/")
		// TODO: find more elegant solution for handling leading slash case
		if len(ss[0]) == 0 {
			ss = ss[1:]
		}
		if len(ss) != 4 {
			log.Println("Bad request", r.RequestURI)
			http.Error(w, "URI should be update/counter/<NAME>/<VALUE>", http.StatusBadRequest)
			return
		}
		name := ss[2]
		svalue := ss[3]
		if len(name) == 0 {
			http.Error(w, "NAME must not be empty", http.StatusBadRequest)
			return
		}
		value, err := strconv.ParseInt(svalue, 10, 64)
		if err != nil {
			http.Error(w, "VALUE type must be int64", http.StatusBadRequest)
			return
		}
		log.Println(name, value)
		_, err = repo.AddAndSave(name, value)
		if err != nil {
			log.Println("Could not add and save counter value", name, value)
			http.Error(w, "Could not add and save counter value", http.StatusInternalServerError)
			return
		}
	}
}

// Listen update/gauge/*, process update/gauge/<NAME>/<VALUE> (value float64)
func gaugeHandler(repo models.GaugeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO refactor duplicate code
		ss := s.Split(r.URL.Path, "/")
		if len(ss[0]) == 0 {
			ss = ss[1:]
		}
		if len(ss) != 4 {
			log.Println("Bad request", r.RequestURI)
			http.Error(w, "URI should be update/gauge/<NAME>/<VALUE>", http.StatusBadRequest)
			return
		}
		name := ss[2]
		svalue := ss[3]
		if len(name) == 0 {
			http.Error(w, "NAME must not be empty", http.StatusBadRequest)
			return
		}
		value, err := strconv.ParseFloat(svalue, 64)
		if err != nil {
			http.Error(w, "VALUE type must be float64", http.StatusBadRequest)
			return
		}
		log.Println(name, value)
		_, err = repo.Save(name, value)
		if err != nil {
			log.Println("Could not save gauge value", name, value)
			http.Error(w, "Could not save gauge value", http.StatusInternalServerError)
			return
		}
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(r.Method, r.RequestURI, err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println(r.Method, r.RequestURI, string(b))
	http.NotFound(w, r)
}
