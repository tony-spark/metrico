package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSingleValueRepository(t *testing.T) {
	r := NewSingleValueRepository()

	t.Run("gauge not found", func(t *testing.T) {
		gauge, err := r.GetGaugeByName(context.Background(), "test")
		assert.Nil(t, gauge)
		assert.Nil(t, err)
	})
	t.Run("gauge saved and found", func(t *testing.T) {
		name := "test1"
		gauge, err := r.SaveGauge(context.Background(), name, float64(3.14))
		assert.NotNil(t, gauge)
		assert.Nil(t, err)
		gauge, err = r.GetGaugeByName(context.Background(), name)
		assert.NotNil(t, gauge)
		assert.Nil(t, err)
	})
	t.Run("gauge value", func(t *testing.T) {
		name := "test2"
		value := float64(3.15)
		gauge1, err := r.SaveGauge(context.Background(), name, value)
		assert.NotNil(t, gauge1)
		assert.Nil(t, err)
		gauge2, err := r.GetGaugeByName(context.Background(), name)
		assert.NotNil(t, gauge2)
		assert.Nil(t, err)
		assert.Equal(t, gauge2.Value, value)
	})
	t.Run("counter not found", func(t *testing.T) {
		counter, err := r.GetCounterByName(context.Background(), "test")
		assert.Nil(t, counter)
		assert.Nil(t, err)
	})
	t.Run("counter saved and found", func(t *testing.T) {
		name := "test1"
		counter, err := r.AddAndSaveCounter(context.Background(), name, int64(314))
		assert.NotNil(t, counter)
		assert.Nil(t, err)
		counter, err = r.GetCounterByName(context.Background(), name)
		assert.NotNil(t, counter)
		assert.Nil(t, err)
	})
	t.Run("counter value", func(t *testing.T) {
		name := "test2"
		values := []int64{1, 4, 5}
		sums := []int64{1, 5, 10}
		for i := 0; i < len(values); i++ {
			counter, err := r.AddAndSaveCounter(context.Background(), name, values[i])
			assert.NotNil(t, counter)
			assert.Nil(t, err)
			assert.Equal(t, counter.Value, sums[i])
		}
	})
}
