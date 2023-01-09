package metrics

import (
	"testing"

	"github.com/rs/zerolog"
)

func BenchmarkAllCollectors(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	cs := []MetricCollector{
		NewMemoryMetricCollector(),
		NewPsUtilMetricsCollector(),
		NewRandomMetricCollector(),
	}
	b.Run("all collectors update ang get metrics data", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, mc := range cs {
				mc.Update()
				for _, m := range mc.Metrics() {
					_ = m.ID()
					_ = m.Type()
					_ = m.String()
					_ = m.Val()
				}
			}
		}
	})
}
