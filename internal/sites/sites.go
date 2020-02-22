package sites

import (
	"net/url"
	"time"
)

type Site struct {
	Name    string
	Url     *url.URL
	Alive   bool
	Latency time.Duration
}

func (s *Site) MarkAvailable(latency time.Duration) {
	s.Alive = true
	s.Latency = latency
}

func (s *Site) MarkUnavailable() {
	s.Alive = false
}
