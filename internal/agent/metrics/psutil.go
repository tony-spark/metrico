package metrics

import (
	"fmt"
	"sync"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/tony-spark/metrico/internal/model"
)

type PsUtilMetricsCollector struct {
	ms       []model.Metric
	vm       *mem.VirtualMemoryStat
	cpuloads []float64
	mu       sync.RWMutex
}

type PsUtilMetric struct {
	name    string
	valueFn func() float64
}

func NewPsUtilMetricsCollector() *PsUtilMetricsCollector {
	psc := &PsUtilMetricsCollector{}

	psc.ms = append(psc.ms, &PsUtilMetric{
		name:    "TotalMemory",
		valueFn: psc.totalMemory,
	}, &PsUtilMetric{
		name:    "FreeMemory",
		valueFn: psc.freeMemory,
	})

	cpus, _ := cpu.Counts(true)
	for i := 0; i < cpus; i++ {
		psc.ms = append(psc.ms, &PsUtilMetric{
			name:    fmt.Sprintf("CPUutilization%d", i+1),
			valueFn: psc.cpuLoad(i),
		})
	}

	return psc
}

func (p *PsUtilMetricsCollector) Metrics() []model.Metric {
	return p.ms
}

func (p *PsUtilMetricsCollector) Update() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.vm, _ = mem.VirtualMemory()
	p.cpuloads, _ = cpu.Percent(0, true)
}

func (p *PsUtilMetricsCollector) totalMemory() float64 {
	if p.vm == nil {
		return 0
	}
	p.mu.RLock()
	defer p.mu.RUnlock()

	return float64(p.vm.Total)
}

func (p *PsUtilMetricsCollector) freeMemory() float64 {
	if p.vm == nil {
		return 0
	}
	p.mu.RLock()
	defer p.mu.RUnlock()

	return float64(p.vm.Free)
}

func (p *PsUtilMetricsCollector) cpuLoad(i int) func() float64 {
	return func() float64 {
		p.mu.RLock()
		defer p.mu.RUnlock()

		return p.cpuloads[i]
	}
}

func (p PsUtilMetric) String() string {
	return fmt.Sprintf("%f", p.valueFn())
}

func (p PsUtilMetric) ID() string {
	return p.name
}

func (p PsUtilMetric) Type() string {
	return model.GAUGE
}

func (p PsUtilMetric) Val() interface{} {
	return p.valueFn()
}
