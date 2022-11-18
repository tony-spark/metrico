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
		rBefore := NewSingleValueRepository()
		_, err = rBefore.SaveGauge(context.Background(), "TestGauge", 1.0)
		assert.Nil(t, err)
		_, err = rBefore.SaveCounter(context.Background(), "TestCounter", 13)
		assert.Nil(t, err)
		err = jfp.Save(context.Background(), rBefore)
		assert.Nil(t, err)
		rAfter := NewSingleValueRepository()
		err = jfp.Load(context.Background(), rAfter)
		assert.Nil(t, err)
		assert.Equal(t, rBefore, rAfter)
	})
}
