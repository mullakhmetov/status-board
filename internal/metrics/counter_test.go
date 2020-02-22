package metrics

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStandardCounter(t *testing.T) {
	c := NewStandardCounter("some name")

	assert.Equal(t, c.Name(), "some name checks")

	wg := sync.WaitGroup{}

	count := 5
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			c.Inc()
			wg.Done()
		}(&wg)
	}
	wg.Wait()

	assert.Equal(t, c.Count(), int64(count))
}
