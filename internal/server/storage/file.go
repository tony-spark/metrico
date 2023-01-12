package storage

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/tony-spark/metrico/internal/model"
	"github.com/tony-spark/metrico/internal/server/models"
)

type JSONFilePersistence struct {
	file *os.File
}

type data struct {
	Gauges   []models.GaugeValue
	Counters []models.CounterValue
}

func (fp JSONFilePersistence) Load(ctx context.Context, r models.MetricRepository) error {
	log.Printf("Loading from %v", fp.file.Name())
	_, err := fp.file.Seek(0, 0)
	if err != nil {
		return err
	}
	bs, err := io.ReadAll(fp.file)
	if len(bs) == 0 && err == nil {
		return nil
	}
	if err != nil {
		return err
	}
	var d data
	err = json.Unmarshal(bs, &d)
	if err != nil {
		return err
	}
	for _, g := range d.Gauges {
		log.Debug().Msgf("Loaded gauge %v = %v", g.Name, g.Value)
		_, err := r.SaveGauge(ctx, g.Name, g.Value)
		if err != nil {
			log.Error().Err(err).Msg("error saving gauge to repository")
		}
	}
	for _, c := range d.Counters {
		log.Debug().Msgf("Loaded counter %v = %v", c.Name, c.Value)
		_, err := r.SaveCounter(ctx, c.Name, c.Value)
		if err != nil {
			log.Error().Err(err).Msg("error saving counter to repository")
		}
	}
	return nil
}

func (fp JSONFilePersistence) Save(ctx context.Context, r models.MetricRepository) error {
	// TODO make save operation atomic
	log.Debug().Msgf("Saving metrics to %v", fp.file.Name())
	ms, err := r.GetAll(ctx)
	if err != nil {
		return err
	}
	gauges := make([]models.GaugeValue, 0)
	counters := make([]models.CounterValue, 0)

	for _, m := range ms {
		switch m.Type() {
		case model.GAUGE:
			gauges = append(gauges, models.GaugeValue{
				Name:  m.ID(),
				Value: m.Val().(float64),
			})
		case model.COUNTER:
			counters = append(counters, models.CounterValue{
				Name:  m.ID(),
				Value: m.Val().(int64),
			})
		}
	}

	d := data{
		Gauges:   gauges,
		Counters: counters,
	}
	bs, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	_, err = fp.file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = fp.file.Write(bs)
	if err != nil {
		return err
	}
	return nil
}

func (fp JSONFilePersistence) Close() error {
	return fp.file.Close()
}

func NewJSONFilePersistence(filename string) (*JSONFilePersistence, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	return &JSONFilePersistence{
		file: f,
	}, nil
}
