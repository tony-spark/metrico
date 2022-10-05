package transports

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tony-spark/metrico/internal"
	"github.com/tony-spark/metrico/internal/agent/metrics"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPTransportGauge(t *testing.T) {
	name := "TestGauge"
	typeName := internal.GAUGE
	value := float64(1.001)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Run("Gauge_validURL", func(t *testing.T) {
			assert.Equal(t, fmt.Sprintf("/update/%v/%v/%v", typeName, name, value), r.URL.Path)
		})
	}))
	defer server.Close()

	transport := NewHTTPTransport(server.URL)
	err := transport.SendMetric(metrics.NewGaugeMetric(name, value))
	t.Run("Gauge_errIsNil", func(t *testing.T) {
		assert.Nil(t, err)
	})
}

func TestHTTPTransportCounter(t *testing.T) {
	name := "TestCounter"
	typeName := internal.COUNTER
	value := int64(12345)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Run("Counter_validURL", func(t *testing.T) {
			assert.Equal(t, fmt.Sprintf("/update/%v/%v/%v", typeName, name, value), r.URL.Path)
		})
	}))
	defer server.Close()

	transport := NewHTTPTransport(server.URL)
	err := transport.SendMetric(metrics.NewCounterMetric(name, value))
	t.Run("Counter_errIsNil", func(t *testing.T) {
		assert.Nil(t, err)
	})
}

func TestHTTPTransportBadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Test bad status", http.StatusBadRequest)
	}))
	defer server.Close()

	transport := NewHTTPTransport(server.URL)
	err := transport.SendMetric(metrics.NewCounterMetric("Test", 0))
	t.Run("Counter_errNotNil", func(t *testing.T) {
		assert.NotNil(t, err)
	})
}

// TODO: rework this check to be safer? (e.g. use mock transport)
func TestHTTPTransportConnectionProblem(t *testing.T) {
	transport := NewHTTPTransport("http://doesnotexists:1010")
	err := transport.SendMetric(metrics.NewCounterMetric("Test", 0))
	t.Run("Counter_errNotNil", func(t *testing.T) {
		assert.NotNil(t, err)
	})
}
