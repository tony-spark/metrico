package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/tony-spark/metrico/internal"
	"github.com/tony-spark/metrico/internal/server/models"
	"html/template"
	"log"
	"net/http"
	"sort"
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
