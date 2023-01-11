package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/model"
	"github.com/tony-spark/metrico/internal/server/models"
)

func checkContentType(w http.ResponseWriter, r *http.Request) error {
	ctype := r.Header.Get("Content-Type")
	t, _, err := mime.ParseMediaType(ctype)
	if err != nil || t != "application/json" {
		http.Error(w, "Only application/json supported", http.StatusUnsupportedMediaType)
		return err
	}
	return nil
}

func readMetric(w http.ResponseWriter, r *http.Request) (*dto.Metric, error) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read body", http.StatusBadRequest)
		return nil, err
	}
	var m dto.Metric
	err = json.Unmarshal(body, &m)
	if err != nil {
		http.Error(w, "Could not parse json", http.StatusBadRequest)
		return nil, err
	}
	if m.MType != model.GAUGE && m.MType != model.COUNTER {
		http.Error(w, "Unknown metric type", http.StatusBadRequest)
		return nil, fmt.Errorf("unknown metric type: %v", m.MType)
	}
	return &m, nil
}

func readMetrics(w http.ResponseWriter, r *http.Request) ([]dto.Metric, error) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read body", http.StatusBadRequest)
		return nil, err
	}
	var ms []dto.Metric
	err = json.Unmarshal(body, &ms)
	if err != nil {
		http.Error(w, "Could not parse json", http.StatusBadRequest)
		return nil, err
	}
	for _, m := range ms {
		if m.MType != model.GAUGE && m.MType != model.COUNTER {
			http.Error(w, "Unknown metric type", http.StatusBadRequest)
			return nil, fmt.Errorf("unknown metric type: %v", m.MType)
		}
	}
	return ms, nil
}

func (router Router) checkHash(mdto dto.Metric, w http.ResponseWriter) bool {
	if router.h != nil {
		ok, err := router.h.Check(mdto)
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
			http.Error(w, "could not check metric integrity", http.StatusInternalServerError)
			return false
		}
		if !ok {
			http.Error(w, "metric integrity check failed (wrong hash?)", http.StatusBadRequest)
		}
		return ok
	}
	return true
}

// UpdatePostHandler godoc
// @Summary Update metric value
// @Accepts json
// @Produce json
// @Param metric_data body dto.Metric true "Metric's data"
// @Success 200 {object} dto.Metric
// @Router /update [post]
func (router Router) UpdatePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := checkContentType(w, r); err != nil {
			log.Error().Err(err).Msg("Wrong content type")
			return
		}
		mdto, err := readMetric(w, r)
		if err != nil {
			log.Error().Err(err).Msg("Could not parse metric")
			return
		}
		if !router.checkHash(*mdto, w) {
			return
		}
		if !mdto.HasValue() {
			http.Error(w, "metric value is null", http.StatusBadRequest)
			return
		}
		mvalue := models.FromDTO(*mdto)
		updated, err := router.ms.UpdateMetric(context.Background(), mvalue)
		if err != nil {
			log.Error().Err(err).Msg("could not save metric")
			http.Error(w, "could not save metric", http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(dto.NewMetric(updated))
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

// BulkUpdatePostHandler godoc
// @Summary Update metric value of multiple metrics
// @Accepts json
// @Produce json
// @Param metric_data body []dto.Metric true "Metric's data"
// @Router /updates [post]
func (router Router) BulkUpdatePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := checkContentType(w, r); err != nil {
			log.Error().Err(err).Msg("Wrong content type")
			return
		}
		ms, err := readMetrics(w, r)
		if err != nil {
			return
		}
		for _, m := range ms {
			if !router.checkHash(m, w) {
				http.Error(w, "hash check failed", http.StatusBadRequest)
				return
			}
		}
		gs := make([]models.GaugeValue, 0)
		cs := make([]models.CounterValue, 0)
		for _, m := range ms {
			if !m.HasValue() {
				http.Error(w, "metric value is null", http.StatusBadRequest)
				return
			}
			switch m.MType {
			case model.GAUGE:
				gs = append(gs, models.GaugeValue{
					Name:  m.ID,
					Value: *m.Value,
				})
			case model.COUNTER:
				cs = append(cs, models.CounterValue{
					Name:  m.ID,
					Value: *m.Delta,
				})
			}
		}
		err = router.ms.UpdateAll(context.Background(), gs, cs)
		if err != nil {
			log.Error().Err(err).Msg("Error saving metrics")
			http.Error(w, "Could not save metrics", http.StatusInternalServerError)
		}
		w.Write([]byte(""))
	}
}

// GetPostHandler godoc
// @Summary Get metric value
// @Accepts json
// @Produce json
// @Param metric_data body dto.Metric true "Metric's data"
// @Success 200 {object} dto.Metric
// @Router /value [post]
func (router Router) GetPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := checkContentType(w, r); err != nil {
			log.Error().Err(err).Msg("Wrong content type")
			return
		}
		mdto, err := readMetric(w, r)
		if err != nil {
			log.Error().Err(err).Msg("Could not parse metric")
			return
		}
		mvalue, err := router.ms.Get(context.Background(), mdto.ID, mdto.MType)
		if err != nil {
			log.Error().Err(err).Msg("Could not get metric")
			http.Error(w, "could not retrieve metric", http.StatusInternalServerError)
			return
		}
		if mvalue == nil {
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}
		mdto = dto.NewMetric(mvalue)
		if router.h != nil {
			mdto.Hash, err = router.h.Hash(*mdto)
			if err != nil {
				http.Error(w, "could not calculate hash for integrity", http.StatusInternalServerError)
				return
			}
		}
		b, err := json.Marshal(mdto)
		if err != nil {
			log.Error().Err(err).Msg("error unmarshalling")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

// MetricGetHandler godoc
// @Summary Get metric value
// @Param metric_type path string true "Metric type" Enum(gauge, counter)
// @Param metric_name path string true "Metric name"
// @Success 200 {string} string "Metric value"
// @Router /value/{metric_type}/{metric_name} [post]
func (router Router) MetricGetHandler(mType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		m, err := router.ms.Get(context.Background(), name, mType)
		if err != nil {
			log.Error().Err(err).Msg("error getting value")
			http.Error(w, "error retrieving value", http.StatusInternalServerError)
			return
		}
		if m == nil {
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprint(m.Val())))
	}
}

// CounterPostHandler godoc
// @Summary Update counter value
// @Param metric_name path string true "Counter name"
// @Param metric_value path int true "Counter value"
// @Router /update/counter/{metric_name}/{metric_value} [post]
func (router Router) CounterPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		svalue := chi.URLParam(r, "svalue")

		value, err := strconv.ParseInt(svalue, 10, 64)
		if err != nil {
			http.Error(w, "VALUE type must be int64", http.StatusBadRequest)
			return
		}
		_, err = router.ms.UpdateCounter(context.Background(), models.CounterValue{Name: name, Value: value})
		if err != nil {
			log.Error().Err(err).Msgf("Could not add and save counter value %s = %v", name, value)
			http.Error(w, "Could not add and save counter value", http.StatusInternalServerError)
			return
		}
	}
}

// GaugePostHandler godoc
// @Summary Update gauge value
// @Param metric_name path string true "Gauge name"
// @Param metric_value path number true "Gauge value"
// @Router /update/gauge/{metric_name}/{metric_value} [post]
func (router Router) GaugePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		svalue := chi.URLParam(r, "svalue")

		value, err := strconv.ParseFloat(svalue, 64)
		if err != nil {
			http.Error(w, "VALUE type must be float64", http.StatusBadRequest)
			return
		}
		_, err = router.ms.UpdateGauge(context.Background(), models.GaugeValue{Name: name, Value: value})
		if err != nil {
			log.Error().Err(err).Msgf("Could not save gauge value %s = %v", name, value)
			http.Error(w, "Could not save gauge value", http.StatusInternalServerError)
			return
		}
	}
}

func (router Router) MetricsViewPageHandler() http.HandlerFunc {
	type Item struct {
		Name  string
		Type  string
		Value string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Items []Item
		}{}

		ms, err := router.ms.GetAll(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("error getting metrics")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, m := range ms {
			data.Items = append(data.Items, Item{m.ID(), m.Type(), fmt.Sprint(m.Val())})
		}

		sort.Slice(data.Items, func(i, j int) bool {
			return data.Items[i].Name < data.Items[j].Name
		})

		w.Header().Set("Content-Type", "text/html; charset=UTF-8")

		err = router.templates.MetricsViewTemplate().Execute(w, data)
		if err != nil {
			log.Error().Err(err).Msg("Error rendering webpage")
			http.Error(w, "Could not display metrics", http.StatusInternalServerError)
			return
		}
	}
}

// PingHandler godoc
// @Summary Get database connection status
// @Success 200
// @Failure 503 {string} string "DB connection is not configured"
// @Failure 500 {string} string "could not check DB or DB is not OK"
// @Router /ping [get]
func (router Router) PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if router.dbm == nil {
			http.Error(w, "DB connection is not configured", http.StatusServiceUnavailable)
			return
		}
		ok, err := router.dbm.Check(context.Background())
		if err != nil || !ok {
			http.Error(w, "could not check DB or DB is not OK", http.StatusInternalServerError)
			return
		}
	}
}

func handleUnknown(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "*")
	http.Error(w, "unknown metric type in "+mtype, http.StatusNotImplemented)
}
