package storage

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestJSONFilePersistence(t *testing.T) {
	t.Run("simple save and load", func(t *testing.T) {
		tempf, err := os.CreateTemp(os.TempDir(), "json-persistence-test")
		if err != nil {
			defer tempf.Close()
		}
		assert.Nil(t, err)
		jfp, err := NewJSONFilePersistence(tempf.Name())
		assert.Nil(t, err)
		grBefore := NewSingleValueGaugeRepository()
		crBefore := NewSingleValueCounterRepository()
		_, err = grBefore.Save("TestGauge", 1.0)
		assert.Nil(t, err)
		_, err = crBefore.Save("TestCounter", 13)
		assert.Nil(t, err)
		err = jfp.Save(grBefore, crBefore)
		assert.Nil(t, err)
		grAfter := NewSingleValueGaugeRepository()
		crAfter := NewSingleValueCounterRepository()
		err = jfp.Load(grAfter, crAfter)
		assert.Nil(t, err)
		assert.Equal(t, grBefore, grAfter)
		assert.Equal(t, crBefore, crAfter)
	})
}
