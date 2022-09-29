package handlers

import (
	"errors"
	"github.com/tony-spark/metrico/internal/server/models"
	"io"
	"log"
	"net/http"
	"strconv"
	s "strings"
)

func check(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		err := errors.New("only POST supported")
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return err
	}
	// TODO temporarily disabled
	//if r.Header.Get("Content-Type") != "text/plain" {
	//	err := errors.New("only text/plain supported")
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return err
	//}
	return nil
}

func extractNameValue(w http.ResponseWriter, r *http.Request) (name string, value string, err error) {
	ss := s.Split(s.Trim(r.URL.Path, "/"), "/")
	if len(ss) == 2 {
		err = errors.New("NAME and VALUE are empty")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if len(ss) < 4 {
		err = errors.New("path should be update/counter/<NAME>/<VALUE>, got " + r.RequestURI)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(ss) > 4 {
		err = errors.New("only path supported: update/gauge/<NAME>/<VALUE>")
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}
	name = ss[2]
	value = ss[3]
	if len(name) == 0 {
		err = errors.New("NAME must not be empty")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	return
}

// CounterHandler listens update/counter/*, processes update/counter/<NAME>/<VALUE> (value int64)
func CounterHandler(repo models.CounterRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := check(w, r)
		if err != nil {
			log.Println(err.Error())
			return
		}
		name, svalue, err := extractNameValue(w, r)
		if err != nil {
			log.Println(err.Error())
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

// GaugeHandler listens update/gauge/*, processes update/gauge/<NAME>/<VALUE> (value float64)
func GaugeHandler(repo models.GaugeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := check(w, r)
		if err != nil {
			log.Println(err.Error())
			return
		}
		name, svalue, err := extractNameValue(w, r)
		if err != nil {
			log.Println(err.Error())
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

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(r.Method, r.RequestURI, err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println(r.Method, r.RequestURI, string(b))
	http.NotFound(w, r)
}
