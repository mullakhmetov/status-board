// Package asker provides some abstract resources availability check functionality.
// It defines Service interface and implements asker for http services.

package asker

import (
	"context"
	"fmt"
	"time"
)

type NotFoundError struct {
	siteName string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Unknown site: %s", e.siteName)
}

type NoResponse struct{}

func (e *NoResponse) Error() string {
	return "No sites"
}

// Response represents resource availability status
type Response struct {
	Name    string
	Alive   bool
	Latency time.Duration
}

// Service defines interface to check resources availability
type Service interface {
	Run(ctx context.Context)
	CheckAll(ctx context.Context) error

	Get(ctx context.Context, name string) (Response, error)
	GetMin(ctx context.Context) (Response, error)
	GetMax(ctx context.Context) (Response, error)
	GetRandom(ctx context.Context) (Response, error)

	Close()
}
