package metrics

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMetric(t *testing.T) {
	router, registry := setupRouter()
	registry.AddCounter("foo")
	registry.Counters["foo"].Inc()

	registry.AddCounter("bar")
	registry.Counters["bar"].Inc()
	registry.Counters["bar"].Inc()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metrics/foo", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, response["name"], "foo checks")
	assert.Equal(t, response["value"].(float64), float64(1))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/metrics/bar", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, response["name"], "bar checks")
	assert.Equal(t, response["value"].(float64), float64(2))
}

func TestMetrics(t *testing.T) {
	router, registry := setupRouter()
	registry.AddCounter("foo")
	registry.Counters["foo"].Inc()

	registry.AddCounter("bar")
	registry.Counters["bar"].Inc()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metrics", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), response["foo checks"].(float64))
	assert.Equal(t, float64(1), response["bar checks"].(float64))
}

func TestMetric_NotFound(t *testing.T) {
	router, registry := setupRouter()
	registry.AddCounter("foo")
	registry.AddCounter("bar")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metrics/buz", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func setupRouter() (*gin.Engine, *Registry) {
	r := gin.Default()
	registry := NewRegistry(false)
	RegisterHandlers(r, registry)
	return r, registry
}
