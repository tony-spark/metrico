package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestJSONFilePersistence(t *testing.T) {
	t.Run("simple save and load", func(t *testing.T) {
		tempf, err := os.CreateTemp(os.TempDir(), "json-persistence-test")
		if err != nil {
			defer tempf.Close()
		}
		require.Nil(t, err)
		jfp, err := NewJSONFilePersistence(tempf.Name())
		require.Nil(t, err)
		grBefore := NewSingleValueGaugeRepository()
		crBefore := NewSingleValueCounterRepository()
		_, err = grBefore.Save(context.Background(), "TestGauge", 1.0)
		assert.Nil(t, err)
		_, err = crBefore.Save(context.Background(), "TestCounter", 13)
		assert.Nil(t, err)
		err = jfp.Save(context.Background(), grBefore, crBefore)
		assert.Nil(t, err)
		grAfter := NewSingleValueGaugeRepository()
		crAfter := NewSingleValueCounterRepository()
		err = jfp.Load(context.Background(), grAfter, crAfter)
		assert.Nil(t, err)
		assert.Equal(t, grBefore, grAfter)
		assert.Equal(t, crBefore, crAfter)
	})
}
