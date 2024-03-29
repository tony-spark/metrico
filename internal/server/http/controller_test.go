package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tony-spark/metrico/internal/server/services"

	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/model"
	"github.com/tony-spark/metrico/internal/server/storage"
)

func TestRouter(t *testing.T) {
	mr := storage.NewSingleValueRepository()
	r := NewController(services.NewMetricService(mr, nil))
	ts := httptest.NewServer(r.r)
	defer ts.Close()

	t.Run("metrics page", func(t *testing.T) {
		statusCode, body := testRequest(t, ts, "GET", "/")
		assert.Equal(t, http.StatusOK, statusCode)
		assert.True(t, strings.Contains(body, "<body>"))
	})
	t.Run("unknown metric type", func(t *testing.T) {
		statusCode, _ := testRequest(t, ts, "POST", "/update/unknown/testCounter/100")
		assert.Equal(t, http.StatusNotImplemented, statusCode)
	})
	t.Run("empty gauge path", func(t *testing.T) {
		statusCode, _ := testRequest(t, ts, "POST", "/update/gauge/")
		assert.Equal(t, http.StatusNotFound, statusCode)
	})
	t.Run("empty counter path", func(t *testing.T) {
		statusCode, _ := testRequest(t, ts, "POST", "/update/counter/")
		assert.Equal(t, http.StatusNotFound, statusCode)
	})
	t.Run("invalid gauge value", func(t *testing.T) {
		statusCode, _ := testRequest(t, ts, "POST", "/update/gauge/wrong/a1.05")
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})
	t.Run("invalid counter value", func(t *testing.T) {
		statusCode, _ := testRequest(t, ts, "POST", "/update/counter/wrong/1.05")
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})
	t.Run("valid counter", func(t *testing.T) {
		statusCode, _ := testRequest(t, ts, "POST", "/update/counter/ok/105")
		assert.Equal(t, http.StatusOK, statusCode)
	})
	t.Run("valid gauge", func(t *testing.T) {
		statusCode, _ := testRequest(t, ts, "POST", "/update/gauge/ok/-13.4523")
		assert.Equal(t, http.StatusOK, statusCode)
	})
	t.Run("wrong method on counter", func(t *testing.T) {
		statusCode, _ := testRequest(t, ts, "GET", "/update/counter/test/105")
		assert.Equal(t, http.StatusMethodNotAllowed, statusCode)
	})
	t.Run("wrong method on gauge", func(t *testing.T) {
		statusCode, _ := testRequest(t, ts, "GET", "/update/gauge/test/-13.4523")
		assert.Equal(t, http.StatusMethodNotAllowed, statusCode)
	})
	t.Run("test gauge not found", func(t *testing.T) {
		statusCode, _ := testRequest(t, ts, "GET", "/value/gauge/absent")
		assert.Equal(t, http.StatusNotFound, statusCode)
	})
	t.Run("test counter not found", func(t *testing.T) {
		statusCode, _ := testRequest(t, ts, "GET", "/value/counter/absent")
		assert.Equal(t, http.StatusNotFound, statusCode)
	})
	t.Run("test gauge read status", func(t *testing.T) {
		statusCode, _ := testRequest(t, ts, "POST", "/update/gauge/test/-13.4523")
		assert.Equal(t, http.StatusOK, statusCode)
		statusCode, _ = testRequest(t, ts, "GET", "/value/gauge/test")
		assert.Equal(t, http.StatusOK, statusCode)
	})
	t.Run("test counter read status", func(t *testing.T) {
		statusCode, _ := testRequest(t, ts, "POST", "/update/counter/test/10")
		assert.Equal(t, http.StatusOK, statusCode)
		statusCode, _ = testRequest(t, ts, "GET", "/value/counter/test")
		assert.Equal(t, http.StatusOK, statusCode)
	})
	t.Run("test gauge value", func(t *testing.T) {
		statusCode, _ := testRequest(t, ts, "POST", "/update/gauge/test1/-12.34")
		assert.Equal(t, http.StatusOK, statusCode)
		statusCode, body := testRequest(t, ts, "GET", "/value/gauge/test1")
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, "-12.34", body)
	})
	t.Run("test counter value", func(t *testing.T) {
		values := []string{"10", "20", "40"}
		sums := []string{"10", "30", "70"}
		for i := 0; i < len(values); i++ {
			statusCode, _ := testRequest(t, ts, "POST", "/update/counter/test1/"+values[i])
			assert.Equal(t, http.StatusOK, statusCode)
			statusCode, body := testRequest(t, ts, "GET", "/value/counter/test1")
			assert.Equal(t, http.StatusOK, statusCode)
			assert.Equal(t, sums[i], body)
		}
	})
	t.Run("test gauge update (post)", func(t *testing.T) {
		v := float64(10.0)
		mreq := &dto.Metric{
			ID:    "UpdateTest1",
			MType: model.GAUGE,
			Value: &v,
		}
		statusCode, mresp := testMetricRequest(t, ts, "POST", "/update/", mreq)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, *mreq.Value, *mresp.Value)

		mvreq := &dto.Metric{
			ID:    mreq.ID,
			MType: mreq.MType,
		}
		statusCode, mvresp := testMetricRequest(t, ts, "POST", "/value", mvreq)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, *mreq.Value, *mvresp.Value)
	})
	t.Run("test counter update (post)", func(t *testing.T) {
		v := int64(10)
		mreq := &dto.Metric{
			ID:    "UpdateTest2",
			MType: model.COUNTER,
			Delta: &v,
		}
		statusCode, mresp := testMetricRequest(t, ts, "POST", "/update/", mreq)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, *mreq.Delta, *mresp.Delta)

		mvreq := &dto.Metric{
			ID:    mreq.ID,
			MType: mreq.MType,
		}
		statusCode, mvresp := testMetricRequest(t, ts, "POST", "/value", mvreq)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, *mreq.Delta, *mvresp.Delta)
	})
	t.Run("test bulk update", func(t *testing.T) {
		gv := 5.5
		dv := int64(4)
		msReq := []dto.Metric{
			{
				ID:    "BulkTestGauge",
				MType: model.GAUGE,
				Value: &gv,
			},
			{
				ID:    "BultTestCounter",
				MType: model.COUNTER,
				Delta: &dv,
			},
		}

		statusCode, _ := testJSONRequest(t, ts, "POST", "/updates", msReq)
		assert.Equal(t, http.StatusOK, statusCode)
	})
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	require.NoError(t, err)

	return resp.StatusCode, string(respBody)
}

func testJSONRequest(t *testing.T, ts *httptest.Server, method, path string, obj interface{}) (int, []byte) {
	var r io.Reader
	if obj != nil {
		r = bytes.NewReader(marshal(t, obj))
	}

	req, err := http.NewRequest(method, ts.URL+path, r)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	require.NoError(t, err)

	return resp.StatusCode, respBody
}

func testMetricRequest(t *testing.T, ts *httptest.Server, method, path string, obj interface{}) (int, *dto.Metric) {
	statusCode, respBody := testJSONRequest(t, ts, method, path, obj)

	var result dto.Metric
	err := json.Unmarshal(respBody, &result)
	require.NoError(t, err)

	return statusCode, &result
}

func marshal(t *testing.T, obj interface{}) []byte {
	b, err := json.Marshal(obj)
	require.NoError(t, err)

	return b
}
