package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tony-spark/metrico/internal/model"
)

func TestRandomMetricCollector(t *testing.T) {
	mc := NewRandomMetricCollector()

	mc.Update()
	t.Run("has one metric", func(t *testing.T) {
		assert.NotEmpty(t, mc.Metrics())
		assert.True(t, len(mc.Metrics()) == 1)
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

func BenchmarkRandomMetricCollector(b *testing.B) {
	mc := NewRandomMetricCollector()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Update()
		mc.Metrics()
	}
}
