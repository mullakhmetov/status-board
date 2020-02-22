package asker

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	router, ms := setupRouter()

	ms.On("Get", mock.AnythingOfType("*gin.Context"), "some-site").Return(Response{}, nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/status/site/some-site", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	ms.AssertExpectations(t)
}

func TestMin(t *testing.T) {
	router, ms := setupRouter()

	ms.On("GetMin", mock.AnythingOfType("*gin.Context")).Return(Response{}, nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/status/min", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	ms.AssertExpectations(t)
}

func TestMax(t *testing.T) {
	router, ms := setupRouter()

	ms.On("GetMax", mock.AnythingOfType("*gin.Context")).Return(Response{}, nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/status/max", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	ms.AssertExpectations(t)
}

func TestRandom(t *testing.T) {
	router, ms := setupRouter()

	ms.On("GetRandom", mock.AnythingOfType("*gin.Context")).Return(Response{}, nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/status/random", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	ms.AssertExpectations(t)
}

func setupRouter() (*gin.Engine, *MockedService) {
	r := gin.Default()
	ms := new(MockedService)
	RegisterHandlers(r, ms)
	return r, ms
}
