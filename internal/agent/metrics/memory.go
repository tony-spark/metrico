package metrics

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/model"
	"runtime"
)

type MemoryMetric struct {
	name          string
	valueFunction func(stats *runtime.MemStats) float64
	collector     *MemoryMetricCollector
}

type PollMetric struct {
	collector *MemoryMetricCollector
}

type MemoryMetricCollector struct {
	metrics      []model.Metric
	memStats     *runtime.MemStats
	refreshCount uint
}

func NewMemoryMetricCollector() *MemoryMetricCollector {
	mmc := &MemoryMetricCollector{
		memStats: &runtime.MemStats{},
	}
	mmc.addCountMetric()
	// TODO: use reflection?
	mmc.addMemoryMetric("Alloc", func(stats *runtime.MemStats) float64 {
		return float64(stats.Alloc)
	})
	mmc.addMemoryMetric("BuckHashSys", func(stats *runtime.MemStats) float64 {
		return float64(stats.BuckHashSys)
	})
	mmc.addMemoryMetric("Frees", func(stats *runtime.MemStats) float64 {
		return float64(stats.Frees)
	})
	mmc.addMemoryMetric("GCCPUFraction", func(stats *runtime.MemStats) float64 {
		return stats.GCCPUFraction
	})
	mmc.addMemoryMetric("GCSys", func(stats *runtime.MemStats) float64 {
		return float64(stats.GCSys)
	})
	mmc.addMemoryMetric("HeapAlloc", func(stats *runtime.MemStats) float64 {
		return float64(stats.HeapAlloc)
	})
	mmc.addMemoryMetric("HeapIdle", func(stats *runtime.MemStats) float64 {
		return float64(stats.HeapIdle)
	})
	mmc.addMemoryMetric("HeapInuse", func(stats *runtime.MemStats) float64 {
		return float64(stats.HeapInuse)
	})
	mmc.addMemoryMetric("HeapObjects", func(stats *runtime.MemStats) float64 {
		return float64(stats.HeapObjects)
	})
	mmc.addMemoryMetric("HeapReleased", func(stats *runtime.MemStats) float64 {
		return float64(stats.HeapReleased)
	})
	mmc.addMemoryMetric("HeapSys", func(stats *runtime.MemStats) float64 {
		return float64(stats.HeapSys)
	})
	mmc.addMemoryMetric("LastGC", func(stats *runtime.MemStats) float64 {
		return float64(stats.LastGC)
	})
	mmc.addMemoryMetric("Lookups", func(stats *runtime.MemStats) float64 {
		return float64(stats.Lookups)
	})
	mmc.addMemoryMetric("MCacheInuse", func(stats *runtime.MemStats) float64 {
		return float64(stats.MCacheInuse)
	})
	mmc.addMemoryMetric("MCacheSys", func(stats *runtime.MemStats) float64 {
		return float64(stats.MCacheSys)
	})
	mmc.addMemoryMetric("MSpanInuse", func(stats *runtime.MemStats) float64 {
		return float64(stats.MSpanInuse)
	})
	mmc.addMemoryMetric("MSpanSys", func(stats *runtime.MemStats) float64 {
		return float64(stats.MSpanSys)
	})
	mmc.addMemoryMetric("Mallocs", func(stats *runtime.MemStats) float64 {
		return float64(stats.Mallocs)
	})
	mmc.addMemoryMetric("NextGC", func(stats *runtime.MemStats) float64 {
		return float64(stats.NextGC)
	})
	mmc.addMemoryMetric("NumForcedGC", func(stats *runtime.MemStats) float64 {
		return float64(stats.NumForcedGC)
	})
	mmc.addMemoryMetric("NumGC", func(stats *runtime.MemStats) float64 {
		return float64(stats.NumGC)
	})
	mmc.addMemoryMetric("OtherSys", func(stats *runtime.MemStats) float64 {
		return float64(stats.OtherSys)
	})
	mmc.addMemoryMetric("PauseTotalNs", func(stats *runtime.MemStats) float64 {
		return float64(stats.PauseTotalNs)
	})
	mmc.addMemoryMetric("StackInuse", func(stats *runtime.MemStats) float64 {
		return float64(stats.StackInuse)
	})
	mmc.addMemoryMetric("StackSys", func(stats *runtime.MemStats) float64 {
		return float64(stats.StackSys)
	})
	mmc.addMemoryMetric("Sys", func(stats *runtime.MemStats) float64 {
		return float64(stats.Sys)
	})
	mmc.addMemoryMetric("TotalAlloc", func(stats *runtime.MemStats) float64 {
		return float64(stats.TotalAlloc)
	})
	return mmc
}

func (m MemoryMetric) String() string {
	return fmt.Sprint(m.valueFunction(m.collector.memStats))
}

func (m MemoryMetric) ID() string {
	return m.name
}

func (m MemoryMetric) Type() string {
	return model.GAUGE
}

func (m MemoryMetric) Val() interface{} {
	return m.valueFunction(m.collector.memStats)
}

func (p PollMetric) String() string {
	return fmt.Sprint(p.collector.refreshCount)
}

func (p PollMetric) ID() string {
	return "PollCount"
}

func (p PollMetric) Type() string {
	return model.COUNTER
}

func (p PollMetric) Val() interface{} {
	return p.collector.refreshCount
}

func (c *MemoryMetricCollector) Update() {
	log.Trace().Msg("Reading memory statistics")
	runtime.ReadMemStats(c.memStats)
	c.refreshCount++
}

func (c *MemoryMetricCollector) Metrics() []model.Metric {
	return c.metrics
}

func (c *MemoryMetricCollector) addMemoryMetric(name string, valueFunction func(stats *runtime.MemStats) float64) {
	metric := MemoryMetric{
		name:          name,
		valueFunction: valueFunction,
		collector:     c,
	}
	c.metrics = append(c.metrics, metric)
}

func (c *MemoryMetricCollector) addCountMetric() {
	metric := PollMetric{
		collector: c,
	}
	c.metrics = append(c.metrics, metric)
}
