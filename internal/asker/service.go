package asker

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/mullakhmetov/status-board/internal/metrics"
	"github.com/mullakhmetov/status-board/internal/sites"
)

type Service interface {
	Run(ctx context.Context)
	CheckAll(ctx context.Context) error

	Get(ctx context.Context, name string) (Response, error)
	GetMin(ctx context.Context) (Response, error)
	GetMax(ctx context.Context) (Response, error)
	GetRandom(ctx context.Context) (Response, error)

	Close()
}

func NewAsker(s sites.Service, metricsRegistry *metrics.Registry, timeout time.Duration, rate time.Duration) Service {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: timeout,
		}).DialContext,
	}
	client := http.Client{Transport: transport}

	// init metric counters
	for _, site := range s.GetAll() {
		metricsRegistry.AddCounter(site.Name)
	}

	return &asker{
		SitesService:    s,
		MetricsRegistry: metricsRegistry,
		httpClient:      client,
		rate:            rate,
	}
}

type asker struct {
	SitesService    sites.Service
	MetricsRegistry *metrics.Registry
	httpClient      http.Client
	rate            time.Duration
}

func (a *asker) Run(ctx context.Context) {
	a.CheckAll(ctx)

	ticker := time.NewTicker(a.rate)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case _ = <-ticker.C:
				a.CheckAll(ctx)
			}
		}
	}()
}

func (a *asker) CheckAll(ctx context.Context) error {
	var wg sync.WaitGroup

	for _, site := range a.SitesService.GetAll() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			wg.Add(1)
			go a.checkSite(ctx, site, &wg)
		}
	}
	wg.Wait()

	return nil

}

func (a *asker) Get(ctx context.Context, name string) (r Response, err error) {
	for _, site := range a.SitesService.GetAll() {
		if site.Name == name {
			a.MetricsRegistry.Counters[site.Name].Inc()
			return Response{site.Name, site.Alive, site.Latency}, nil
		}
	}

	return r, &NotFoundError{name}
}

func (a *asker) GetMin(ctx context.Context) (r Response, err error) {
	sorted := a.SitesService.GetSortedByLatency()
	if len(sorted) == 0 {
		return r, &NoResponse{}
	}
	min := sorted[0]

	a.MetricsRegistry.Counters[min.Name].Inc()

	return Response{Name: min.Name, Alive: true, Latency: min.Latency}, nil
}

func (a *asker) GetMax(ctx context.Context) (r Response, err error) {
	sorted := a.SitesService.GetSortedByLatency()
	if len(sorted) == 0 {
		return r, &NoResponse{}
	}

	max := sorted[len(sorted)-1]

	a.MetricsRegistry.Counters[max.Name].Inc()

	return Response{Name: max.Name, Alive: true, Latency: max.Latency}, nil
}

func (a *asker) GetRandom(ctx context.Context) (r Response, err error) {
	sites := a.SitesService.GetAll()
	if len(sites) == 0 {
		return r, &NoResponse{}
	}

	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(sites)
	site := sites[n]

	a.MetricsRegistry.Counters[site.Name].Inc()

	return Response{Name: site.Name, Alive: site.Alive, Latency: site.Latency}, nil
}

// nothing to finalize
func (a *asker) Close() {}

// returns nil if site is available
func (a *asker) checkSite(ctx context.Context, site *sites.Site, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, "GET", site.Url.String(), nil)
	if err != nil {
		log.Printf("[ERROR] failed to make request to %s site: %+v", site.Url.String(), err)
		site.MarkUnavailable()
		return
	}

	start := time.Now()
	resp, err := a.httpClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] request to %s site failed: %+v", site.Url.String(), err)
		site.MarkUnavailable()
		return
	}
	// omit reps body & code
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)

	site.MarkAvailable(time.Since(start))
}
