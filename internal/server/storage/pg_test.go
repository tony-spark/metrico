package storage

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PgTestSuite struct {
	suite.Suite

	pgm *PgDatabaseManager
	gs  []string
	cs  []string
}

func (suite *PgTestSuite) SetupSuite() {
	var err error
	dsn := os.Getenv("TEST_DSN")
	if len(dsn) == 0 {
		dsn = "postgresql://postgres:postgres@localhost:5432/metrico?sslmode=disable"
	}
	suite.pgm, err = NewPgManager(dsn)
	suite.Require().NoError(err)
}

func (suite *PgTestSuite) TestPgStorage() {
	r := suite.pgm.mdb
	suite.Run("ping", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		ok, err := suite.pgm.Check(ctx)
		assert.NoError(suite.T(), err)
		assert.True(suite.T(), ok)
	})
	suite.Run("gauge not found", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		gauge, err := r.GetGaugeByName(ctx, "test")
		assert.Nil(suite.T(), gauge)
		assert.Nil(suite.T(), err)
	})
	suite.Run("gauge saved and found", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		name := "test1"
		suite.gs = append(suite.gs, name)
		gauge, err := r.SaveGauge(ctx, name, float64(3.14))
		assert.NotNil(suite.T(), gauge)
		assert.Nil(suite.T(), err)
		gauge, err = r.GetGaugeByName(ctx, name)
		assert.NotNil(suite.T(), gauge)
		assert.Nil(suite.T(), err)
	})
	suite.Run("gauge value", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		name := "test2"
		suite.gs = append(suite.gs, name)
		value := float64(3.15)
		gauge1, err := r.SaveGauge(ctx, name, value)
		assert.NotNil(suite.T(), gauge1)
		assert.NoError(suite.T(), err)
		gauge2, err := r.GetGaugeByName(ctx, name)
		assert.NotNil(suite.T(), gauge2)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), gauge2.Value, value)
	})
	suite.Run("get all gauges", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		for i := 3; i < 6; i++ {
			name := "test" + fmt.Sprint(i)
			suite.gs = append(suite.gs, name)
			value := 2.71
			_, err := r.SaveGauge(ctx, name, value)
			assert.NoError(suite.T(), err)
		}
		gs, err := r.getAllGauges(ctx)
		assert.NoError(suite.T(), err)
		assert.True(suite.T(), len(gs) >= 3)
	})
	suite.Run("counter not found", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		counter, err := r.GetCounterByName(ctx, "test")
		assert.Nil(suite.T(), counter)
		assert.Nil(suite.T(), err)
	})
	suite.Run("counter saved and found", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		name := "test1"
		suite.cs = append(suite.cs, name)
		counter, err := r.AddAndSaveCounter(ctx, name, int64(314))
		assert.NotNil(suite.T(), counter)
		assert.Nil(suite.T(), err)
		counter, err = r.GetCounterByName(ctx, name)
		assert.NotNil(suite.T(), counter)
		assert.Nil(suite.T(), err)
	})
	suite.Run("counter value", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		name := "test2"
		suite.cs = append(suite.cs, name)
		values := []int64{1, 4, 5}
		sums := []int64{1, 5, 10}
		for i := 0; i < len(values); i++ {
			counter, err := r.AddAndSaveCounter(ctx, name, values[i])
			assert.NotNil(suite.T(), counter)
			assert.Nil(suite.T(), err)
			assert.Equal(suite.T(), counter.Value, sums[i])
		}
	})
	suite.Run("get all counters", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		for i := 3; i < 6; i++ {
			name := "test" + fmt.Sprint(i)
			suite.cs = append(suite.cs, name)
			value := int64(33)
			_, err := r.SaveCounter(ctx, name, value)
			assert.NoError(suite.T(), err)
		}
		gs, err := r.getAllCounters(ctx)
		assert.NoError(suite.T(), err)
		assert.True(suite.T(), len(gs) >= 3)
	})
	suite.Run("get all metrics", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		for i := 6; i < 9; i++ {
			name := "test" + fmt.Sprint(i)
			suite.gs = append(suite.gs, name)
			value := 2.71
			_, err := r.SaveGauge(ctx, name, value)
			assert.NoError(suite.T(), err)
		}
		for i := 6; i < 9; i++ {
			name := "test" + fmt.Sprint(i)
			suite.cs = append(suite.cs, name)
			value := int64(33)
			_, err := r.SaveCounter(ctx, name, value)
			assert.NoError(suite.T(), err)
		}
		ms, err := r.GetAll(ctx)
		assert.NoError(suite.T(), err)
		assert.True(suite.T(), len(ms) >= 6)
	})
}

func (suite *PgTestSuite) TearDownSuite() {
	deleteHelper := func(tableName string) func(string) {
		return func(name string) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := suite.pgm.mdb.deleteMetric(ctx, tableName, name)
			suite.Assert().NoError(err)
			cancel()
		}
	}
	deleteGauge := deleteHelper("gauges")
	deleteCounter := deleteHelper("counters")
	for _, gauge := range suite.gs {
		deleteGauge(gauge)
	}
	for _, counter := range suite.cs {
		deleteCounter(counter)
	}
	err := suite.pgm.Close()
	suite.Require().NoError(err)
}

func TestPgStorage(t *testing.T) {
	suite.Run(t, new(PgTestSuite))
}
