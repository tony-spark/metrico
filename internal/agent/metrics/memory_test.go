package metrics

import (
	"github.com/stretchr/testify/assert"
	"github.com/tony-spark/metrico/internal/model"
	"testing"
)

func TestMemoryMetricCollector(t *testing.T) {
	mc := NewMemoryMetricCollector()

	mc.Update()
	t.Run("has metrics", func(t *testing.T) {
		assert.NotEmpty(t, mc.Metrics())
	})
	t.Run("metrics not empty", func(t *testing.T) {
		for _, m := range mc.Metrics() {
			assert.NotEmpty(t, m.ID())
			assert.NotEmpty(t, m.String())
		}
	})
	t.Run("metrics type", func(t *testing.T) {
		for _, m := range mc.Metrics() {
			switch m.Type() {
			case model.COUNTER:
				_, ok := m.Val().(int64)
				assert.True(t, ok)
			case model.GAUGE:
				_, ok := m.Val().(float64)
				assert.True(t, ok)
			}
		}
	})
}

func BenchmarkMemoryMetricCollector(b *testing.B) {
	mc := NewMemoryMetricCollector()
	for i := 0; i < b.N; i++ {
		mc.Update()
		mc.Metrics()
	}
}
