package storage

import (
	"context"
	"encoding/json"
	"github.com/tony-spark/metrico/internal/server/models"
	"io"
	"log"
	"os"
)

type JSONFilePersistence struct {
	file *os.File
}

type data struct {
	Gauges   []*models.GaugeValue
	Counters []*models.CounterValue
}

func (fp JSONFilePersistence) Load(ctx context.Context, gr models.GaugeRepository, cr models.CounterRepository) error {
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
		gr.Save(ctx, g.Name, g.Value)
		log.Printf("Loaded gauge %v = %v", g.Name, g.Value)
	}
	for _, c := range d.Counters {
		cr.Save(ctx, c.Name, c.Value)
		log.Printf("Loaded counter %v = %v", c.Name, c.Value)
	}
	return nil
}

func (fp JSONFilePersistence) Save(ctx context.Context, gr models.GaugeRepository, cr models.CounterRepository) error {
	// TODO make save operation atomic
	log.Printf("Saving metrics to %v", fp.file.Name())
	gauges, err := gr.GetAll(ctx)
	if err != nil {
		return err
	}
	counters, err := cr.GetAll(ctx)
	if err != nil {
		return err
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
