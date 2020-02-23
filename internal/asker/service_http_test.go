package asker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/mullakhmetov/status-board/internal/metrics"
	"github.com/mullakhmetov/status-board/internal/sites"
	"github.com/stretchr/testify/assert"
)

func TestAsker_NewHttpAsker(t *testing.T) {
	ss := []*sites.Site{
		&sites.Site{Name: "google.com"},
		&sites.Site{Name: "vk.com"},
	}

	mockedSites := new(sites.MockedService)
	mockedSites.On("GetAll").Return(ss)

	mockedMetrics := metrics.Registry{
		InitCounterFunc: func(name string) metrics.Counter {
			mockedCounter := new(metrics.MockedCounter)
			return mockedCounter
		},
		Counters: make(map[string]metrics.Counter),
	}

	_ = NewHttpAsker(mockedSites, &mockedMetrics, time.Second, time.Second)
	assert.Equal(t, len(ss), len(mockedMetrics.Counters))
	mockedSites.AssertExpectations(t)
}

func TestAsker_Run(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "google.com" || r.URL.String() == "vk.com" {
			_, err := w.Write([]byte(""))
			assert.NoError(t, err)
			return
		}
		w.WriteHeader(404)
	}))
	defer ts.Close()

	google, err := url.Parse(ts.URL + "/google.com")
	assert.NoError(t, err)
	vk, err := url.Parse(ts.URL + "/vk.com")
	assert.NoError(t, err)

	ss := []*sites.Site{
		&sites.Site{Name: "google.com", Url: google, Alive: false},
		&sites.Site{Name: "vk.com", Url: vk, Alive: false},
	}

	mockedSites := new(sites.MockedService)
	mockedSites.On("GetAll").Return(ss)

	mockedMetrics := metrics.Registry{
		InitCounterFunc: func(name string) metrics.Counter {
			mockedCounter := new(metrics.MockedCounter)
			return mockedCounter
		},
		Counters: make(map[string]metrics.Counter),
	}

	a := NewHttpAsker(mockedSites, &mockedMetrics, time.Second, time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	a.Run(ctx)

	time.Sleep(time.Second)

	mockedSites.AssertExpectations(t)
	assert.True(t, len(mockedSites.Calls) > int(time.Second/time.Millisecond/2))
}

func TestAsker_Run_Cancel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "google.com" || r.URL.String() == "vk.com" {
			_, err := w.Write([]byte(""))
			assert.NoError(t, err)
			return
		}
		w.WriteHeader(404)
	}))
	defer ts.Close()

	google, err := url.Parse(ts.URL + "/google.com")
	assert.NoError(t, err)
	vk, err := url.Parse(ts.URL + "/vk.com")
	assert.NoError(t, err)

	ss := []*sites.Site{
		&sites.Site{Name: "google.com", Url: google, Alive: false},
		&sites.Site{Name: "vk.com", Url: vk, Alive: false},
	}

	mockedSites := new(sites.MockedService)
	mockedSites.On("GetAll").Return(ss)

	mockedMetrics := metrics.Registry{
		InitCounterFunc: func(name string) metrics.Counter {
			mockedCounter := new(metrics.MockedCounter)
			return mockedCounter
		},
		Counters: make(map[string]metrics.Counter),
	}

	a := NewHttpAsker(mockedSites, &mockedMetrics, time.Second, time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	cancel()
	a.Run(ctx)

	time.Sleep(time.Second)

	mockedSites.AssertExpectations(t)
	mockedSites.AssertNumberOfCalls(t, "GetAll", 2)
}

func TestAsker_CheckAll(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "google.com" || r.URL.String() == "vk.com" {
			_, err := w.Write([]byte(""))
			assert.NoError(t, err)
			return
		}
		w.WriteHeader(404)
	}))
	defer ts.Close()

	google, err := url.Parse(ts.URL + "/google.com")
	assert.NoError(t, err)
	vk, err := url.Parse(ts.URL + "/vk.com")
	assert.NoError(t, err)

	ss := []*sites.Site{
		&sites.Site{Name: "google.com", Url: google, Alive: false},
		&sites.Site{Name: "vk.com", Url: vk, Alive: false},
	}

	mockedSites := new(sites.MockedService)
	mockedSites.On("GetAll").Return(ss)

	mockedMetrics := metrics.Registry{
		InitCounterFunc: func(name string) metrics.Counter {
			mockedCounter := new(metrics.MockedCounter)
			return mockedCounter
		},
		Counters: make(map[string]metrics.Counter),
	}

	a := NewHttpAsker(mockedSites, &mockedMetrics, time.Second, time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.CheckAll(ctx)

	// marked as alive
	for _, s := range ss {
		assert.Equal(t, s.Alive, true)
	}
}

func TestAsker_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "google.com" || r.URL.String() == "vk.com" {
			_, err := w.Write([]byte(""))
			assert.NoError(t, err)
			return
		}
		w.WriteHeader(404)
	}))
	defer ts.Close()

	google, err := url.Parse(ts.URL + "/google.com")
	assert.NoError(t, err)
	vk, err := url.Parse(ts.URL + "/vk.com")
	assert.NoError(t, err)

	ss := []*sites.Site{
		&sites.Site{Name: "google.com", Url: google, Alive: false},
		&sites.Site{Name: "vk.com", Url: vk, Alive: false},
	}

	mockedSites := new(sites.MockedService)
	mockedSites.On("GetAll").Return(ss)

	mockedMetrics := metrics.Registry{
		InitCounterFunc: func(name string) metrics.Counter {
			mockedCounter := new(metrics.MockedCounter)
			return mockedCounter
		},
		Counters: make(map[string]metrics.Counter),
	}

	a := NewHttpAsker(mockedSites, &mockedMetrics, time.Second, time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a.CheckAll(ctx)

	resp, err := a.Get(ctx, "google.com")
	assert.NoError(t, err)
	assert.Equal(t, resp.Alive, true)
	assert.Equal(t, resp.Name, "google.com")

	resp, err = a.Get(ctx, "unknown.site")
	assert.Error(t, err)
	assert.Equal(t, resp.Alive, false)
}

func TestAsker_CheckUnreachable(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	ts.Close()

	url, err := url.Parse(ts.URL)
	assert.NoError(t, err)

	ss := []*sites.Site{
		&sites.Site{Name: "google.com", Url: url, Alive: true},
		&sites.Site{Name: "vk.com", Url: url, Alive: true},
	}

	mockedSites := new(sites.MockedService)
	mockedSites.On("GetAll").Return(ss)

	mockedMetrics := metrics.Registry{
		InitCounterFunc: func(name string) metrics.Counter {
			mockedCounter := new(metrics.MockedCounter)
			return mockedCounter
		},
		Counters: make(map[string]metrics.Counter),
	}

	a := NewHttpAsker(mockedSites, &mockedMetrics, time.Second, time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a.CheckAll(ctx)

	resp, err := a.Get(ctx, "google.com")
	assert.NoError(t, err)
	assert.Equal(t, resp.Alive, false)
	assert.Equal(t, resp.Name, "google.com")

	resp, err = a.Get(ctx, "vk.com")
	assert.NoError(t, err)
	assert.Equal(t, resp.Alive, false)
	assert.Equal(t, resp.Name, "vk.com")

}

func TestAsker_GetMin_GetMax_GetRandom(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "google.com" || r.URL.String() == "vk.com" {
			_, err := w.Write([]byte(""))
			assert.NoError(t, err)
			return
		}
		w.WriteHeader(404)
	}))
	defer ts.Close()

	google, err := url.Parse(ts.URL + "/google.com")
	assert.NoError(t, err)
	vk, err := url.Parse(ts.URL + "/vk.com")
	assert.NoError(t, err)

	ss := []*sites.Site{
		&sites.Site{Name: "google.com", Url: google, Alive: false},
		&sites.Site{Name: "vk.com", Url: vk, Alive: false},
	}

	mockedSites := new(sites.MockedService)
	mockedSites.On("GetAll").Return(ss)
	mockedSites.On("GetSortedByLatency").Return(ss)

	mockedMetrics := metrics.Registry{
		InitCounterFunc: func(name string) metrics.Counter {
			mockedCounter := new(metrics.MockedCounter)
			return mockedCounter
		},
		Counters: make(map[string]metrics.Counter),
	}

	a := NewHttpAsker(mockedSites, &mockedMetrics, time.Second, time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a.CheckAll(ctx)

	resp, err := a.GetMin(ctx)
	assert.NoError(t, err)
	assert.Equal(t, resp.Alive, true)
	assert.Equal(t, resp.Name, "google.com")

	resp, err = a.GetMax(ctx)
	assert.NoError(t, err)
	assert.Equal(t, resp.Alive, true)
	assert.Equal(t, resp.Name, "vk.com")

	resp, err = a.GetRandom(ctx)
	assert.NoError(t, err)
	assert.Equal(t, resp.Alive, true)
}
