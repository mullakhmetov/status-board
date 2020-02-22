package asker

import (
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

type Response struct {
	Name    string
	Alive   bool
	Latency time.Duration
}
