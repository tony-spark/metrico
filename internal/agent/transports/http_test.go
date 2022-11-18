package transports

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tony-spark/metrico/internal/agent/metrics"
	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/hash"
	"github.com/tony-spark/metrico/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPTransportGauge(t *testing.T) {
	name := "TestGauge"
	value := float64(1.001)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Run("gauge valid url and payload", func(t *testing.T) {
			assert.Equal(t, "/update/", r.URL.Path)
			bs, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			var m dto.Metric
			err = json.Unmarshal(bs, &m)
			assert.Nil(t, err)
			expected := dto.Metric{
				ID:    name,
				MType: model.GAUGE,
				Delta: nil,
				Value: &value,
			}
			assert.Equal(t, m, expected)
		})
	}))
	defer server.Close()

	transport := NewHTTPTransport(server.URL)
	err := transport.SendMetric(metrics.NewGaugeMetric(name, value))
	t.Run("send gauge no error", func(t *testing.T) {
		assert.Nil(t, err)
	})
}

func TestHTTPTransportCounter(t *testing.T) {
	name := "TestCounter"
	value := int64(12345)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Run("counter valid payload", func(t *testing.T) {
			assert.Equal(t, "/update/", r.URL.Path)
			defer r.Body.Close()
			bs, err := io.ReadAll(r.Body)
			require.Nil(t, err)
			var m dto.Metric
			err = json.Unmarshal(bs, &m)
			assert.Nil(t, err)
			expected := dto.Metric{
				ID:    name,
				MType: model.COUNTER,
				Delta: &value,
				Value: nil,
			}
			assert.Equal(t, m, expected)
		})
	}))
	defer server.Close()

	transport := NewHTTPTransport(server.URL)
	err := transport.SendMetric(metrics.NewCounterMetric(name, value))
	t.Run("send counter no error", func(t *testing.T) {
		assert.Nil(t, err)
	})
}

func TestHTTPTransportBadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad status", http.StatusBadRequest)
	}))
	defer server.Close()

	transport := NewHTTPTransport(server.URL)
	err := transport.SendMetric(metrics.NewCounterMetric("Test", 0))
	t.Run("counter error", func(t *testing.T) {
		assert.NotNil(t, err)
	})
}

// TODO: rework this check to be safer? (e.g. use mock transport)
func TestHTTPTransportConnectionProblem(t *testing.T) {
	transport := NewHTTPTransport("http://doesnotexist:1010")
	err := transport.SendMetric(metrics.NewCounterMetric("Test", 0))
	t.Run("connection error", func(t *testing.T) {
		assert.NotNil(t, err)
	})
}

func TestHTTPTransportHashed(t *testing.T) {
	h := hash.NewSha256Hmac("key")
	name := "TestCounter"
	value := int64(12345)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Run("counter hash present", func(t *testing.T) {
			assert.Equal(t, "/update/", r.URL.Path)
			defer r.Body.Close()
			bs, err := io.ReadAll(r.Body)
			require.Nil(t, err)
			var m dto.Metric
			err = json.Unmarshal(bs, &m)
			assert.Nil(t, err)
			assert.NotEmpty(t, m.Hash)
			check, err := h.Check(m)
			assert.Nil(t, err)
			assert.True(t, check)
		})
	}))
	defer server.Close()

	transport := NewHTTPTransportHashed(server.URL, h)
	err := transport.SendMetric(metrics.NewCounterMetric(name, value))
	t.Run("send counter with hash no error", func(t *testing.T) {
		assert.Nil(t, err)
	})
}
