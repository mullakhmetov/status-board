package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mullakhmetov/status-board/internal/asker"
	"github.com/mullakhmetov/status-board/internal/metrics"
	"github.com/mullakhmetov/status-board/internal/sites"
)

type services struct {
	sites sites.Service
	asker asker.Service
}

type server struct {
	srv *http.Server
	*services
	terminated chan struct{}
}

type ServerOpts struct {
	Port         int
	Timeout      time.Duration
	ChecksRate   time.Duration
	StoreMetrics bool
	SitesPath    string
}

func NewServer(opts ServerOpts) (*server, error) {
	router := gin.Default()

	sitesServices := sites.NewFileSitesService(opts.SitesPath)
	sitesServices.Warmup()

	metricsRegistry := metrics.NewRegistry(!opts.StoreMetrics)
	metrics.RegisterHandlers(router, metricsRegistry)

	askerService := asker.NewAsker(sitesServices, metricsRegistry, opts.Timeout, opts.ChecksRate)
	asker.RegisterHandlers(router, askerService)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", opts.Port),
		Handler: router,
	}

	s := &server{
		srv: srv,
		services: &services{
			sites: sitesServices,
			asker: askerService,
		},
		terminated: make(chan struct{}),
	}
	return s, nil
}

func (s *server) Run(ctx context.Context) error {
	// start asker loop
	s.services.asker.Run(ctx)

	go func() {
		// Graceful shutdown
		<-ctx.Done()
		// Close services
		s.services.asker.Close()
		s.services.sites.Close()

		s.srv.Shutdown(ctx)
		log.Print("[INFO] server was shut down")
	}()

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	close(s.terminated)
	return nil
}

func (s *server) Wait() {
	<-s.terminated
}
