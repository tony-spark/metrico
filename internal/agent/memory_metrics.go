package agent

import (
	"fmt"
	"github.com/tony-spark/metrico/internal"
	"log"
	"runtime"
)

type MemoryMetric struct {
	name          string
	valueFunction func(stats *runtime.MemStats) string
	collector     *MemoryMetricCollector
}

type PollMetric struct {
	collector *MemoryMetricCollector
}

type MemoryMetricCollector struct {
	metrics      []internal.Metric
	memStats     *runtime.MemStats
	refreshCount uint
}

func NewMemoryMetricCollector() *MemoryMetricCollector {
	mmc := &MemoryMetricCollector{
		memStats: &runtime.MemStats{},
	}
	mmc.addCountMetric()
	// TODO: use reflection?
	mmc.addMemoryMetric("Alloc", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.Alloc)
	})
	mmc.addMemoryMetric("BuckHashSys", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.BuckHashSys)
	})
	mmc.addMemoryMetric("Frees", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.Frees)
	})
	mmc.addMemoryMetric("GCCPUFraction", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.GCCPUFraction)
	})
	mmc.addMemoryMetric("GCSys", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.GCSys)
	})
	mmc.addMemoryMetric("HeapAlloc", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.HeapAlloc)
	})
	mmc.addMemoryMetric("HeapIdle", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.HeapIdle)
	})
	mmc.addMemoryMetric("HeapInuse", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.HeapInuse)
	})
	mmc.addMemoryMetric("HeapObjects", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.HeapObjects)
	})
	mmc.addMemoryMetric("HeapReleased", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.HeapReleased)
	})
	mmc.addMemoryMetric("HeapSys", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.HeapSys)
	})
	mmc.addMemoryMetric("LastGC", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.LastGC)
	})
	mmc.addMemoryMetric("Lookups", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.Lookups)
	})
	mmc.addMemoryMetric("MCacheInuse", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.MCacheInuse)
	})
	mmc.addMemoryMetric("MCacheSys", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.MCacheSys)
	})
	mmc.addMemoryMetric("MSpanInuse", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.MSpanInuse)
	})
	mmc.addMemoryMetric("MSpanSys", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.MSpanSys)
	})
	mmc.addMemoryMetric("Mallocs", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.Mallocs)
	})
	mmc.addMemoryMetric("NextGC", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.NextGC)
	})
	mmc.addMemoryMetric("NumForcedGC", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.NumForcedGC)
	})
	mmc.addMemoryMetric("NumGC", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.NumGC)
	})
	mmc.addMemoryMetric("OtherSys", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.OtherSys)
	})
	mmc.addMemoryMetric("PauseTotalNs", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.PauseTotalNs)
	})
	mmc.addMemoryMetric("StackInuse", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.StackInuse)
	})
	mmc.addMemoryMetric("StackSys", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.StackSys)
	})
	mmc.addMemoryMetric("Sys", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.Sys)
	})
	mmc.addMemoryMetric("TotalAlloc", func(stats *runtime.MemStats) string {
		return fmt.Sprint(stats.TotalAlloc)
	})
	return mmc
}

func (m MemoryMetric) String() string {
	return m.valueFunction(m.collector.memStats)
}

func (m MemoryMetric) Name() string {
	return m.name
}

func (m MemoryMetric) Type() string {
	// TODO: all runtime metrics are gauge?
	return internal.GAUGE
}

func (p PollMetric) String() string {
	return fmt.Sprint(p.collector.refreshCount)
}

func (p PollMetric) Name() string {
	return "PollCount"
}

func (p PollMetric) Type() string {
	return internal.COUNTER
}

func (c *MemoryMetricCollector) Update() {
	// TODO: ReadMemStats causes stopTheWorld, what should we do about it?
	log.Println("Reading memory statistics")
	runtime.ReadMemStats(c.memStats)
	c.refreshCount++
}

func (c *MemoryMetricCollector) Metrics() []internal.Metric {
	return c.metrics
}

func (c *MemoryMetricCollector) addMemoryMetric(name string, valueFunction func(stats *runtime.MemStats) string) {
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
