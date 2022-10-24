package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/tony-spark/metrico/internal"
	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/server/models"
	"html/template"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"sort"
	"strconv"
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

func readMetrics(w http.ResponseWriter, r *http.Request) (*dto.Metrics, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read body", http.StatusBadRequest)
		return nil, err
	}
	var m dto.Metrics
	err = json.Unmarshal(body, &m)
	if err != nil {
		http.Error(w, "Could not parse json", http.StatusBadRequest)
		return nil, err
	}
	if m.MType != internal.GAUGE && m.MType != internal.COUNTER {
		http.Error(w, "Unknown metric type", http.StatusBadRequest)
		return nil, fmt.Errorf("unknown metric type: %v", m.MType)
	}
	return &m, nil
}

func UpdatePostHandler(gaugeRepo models.GaugeRepository, counterRepo models.CounterRepository, postUpdateFn func()) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := checkContentType(w, r); err != nil {
			log.Println(err.Error())
			return
		}
		m, err := readMetrics(w, r)
		if err != nil {
			log.Println(err.Error())
			return
		}
		switch m.MType {
		case internal.GAUGE:
			if m.Value == nil {
				http.Error(w, "gauge value is null", http.StatusBadRequest)
				return
			}
			g, err := gaugeRepo.Save(m.ID, *m.Value)
			if err != nil {
				http.Error(w, "could not save gauge value", http.StatusInternalServerError)
			}
			m.Value = &g.Value
		case internal.COUNTER:
			if m.Delta == nil {
				http.Error(w, "counter value is null", http.StatusBadRequest)
				return
			}
			c, err := counterRepo.AddAndSave(m.ID, *m.Delta)
			if err != nil {
				http.Error(w, "could not update counter value", http.StatusInternalServerError)
			}
			m.Delta = &c.Value
		}
		b, err := json.Marshal(m)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		if postUpdateFn != nil {
			postUpdateFn()
		}
	}
}

func GetPostHandler(gaugeRepo models.GaugeRepository, counterRepo models.CounterRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := checkContentType(w, r); err != nil {
			log.Println(err.Error())
			return
		}
		m, err := readMetrics(w, r)
		if err != nil {
			log.Println(err.Error())
			return
		}
		switch m.MType {
		case internal.GAUGE:
			g, err := gaugeRepo.GetByName(m.ID)
			if err != nil {
				http.Error(w, "could not retrieve gauge value", http.StatusInternalServerError)
				return
			}
			if g == nil {
				http.Error(w, "gauge not found", http.StatusNotFound)
				return
			}
			m.Value = &g.Value
		case internal.COUNTER:
			c, err := counterRepo.GetByName(m.ID)
			if err != nil {
				http.Error(w, "could not retrieve counter value", http.StatusInternalServerError)
				return
			}
			if c == nil {
				http.Error(w, "counter not found", http.StatusNotFound)
				return
			}
			m.Delta = &c.Value
		}
		b, err := json.Marshal(m)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

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

func CounterPostHandler(repo models.CounterRepository, postUpdateFn func()) http.HandlerFunc {
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
		if postUpdateFn != nil {
			postUpdateFn()
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

func GaugePostHandler(repo models.GaugeRepository, postUpdateFn func()) http.HandlerFunc {
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
		if postUpdateFn != nil {
			postUpdateFn()
		}
	}
}

func PageHandler(gaugeRepo models.GaugeRepository, counterRepo models.CounterRepository) http.HandlerFunc {
	const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="refresh" content="10">
		<title>Metrics</title>
	</head>
	<body>
		<table border="1">
			<tr>
				<th>Metric</th>
				<th>Type</th>
				<th>Value</th>
			</tr>
			{{range .Items}}
			<tr>
				<td>{{ .Name }}</td>
				<td>{{ .Type }}</td>
				<td>{{ .Value }}</td>
			</tr>
			{{else}}
			<tr>
				<td colspan="3"><strong>No metrics</strong></td>
			</tr>
			{{end}}
		</table>
	</body>
</html>
`
	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		log.Fatalln("Could not parse template", err)
	}

	type Item struct {
		Name  string
		Type  string
		Value string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Items []Item
		}{}

		// TODO if metrics are being updated during page generation
		gs, err := gaugeRepo.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, g := range gs {
			data.Items = append(data.Items, Item{g.Name, internal.GAUGE, fmt.Sprint(g.Value)})
		}
		vs, err := counterRepo.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, v := range vs {
			data.Items = append(data.Items, Item{v.Name, internal.COUNTER, fmt.Sprint(v.Value)})
		}

		sort.Slice(data.Items, func(i, j int) bool {
			return data.Items[i].Name < data.Items[j].Name
		})

		w.Header().Set("Content-Type", "text/html; charset=UTF-8")

		err = t.Execute(w, data)
		if err != nil {
			http.Error(w, "Could not display metrics", http.StatusInternalServerError)
			return
		}
	}
}
