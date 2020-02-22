package metrics

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	r := NewRegistry(false)

	wg := sync.WaitGroup{}

	expectedStats := make(map[string]int64)
	count := 5
	for i := 0; i < count; i++ {
		wg.Add(1)

		name := "counter" + strconv.Itoa(i)
		expectedStats[name+" checks"] = int64(0)

		go func(*sync.WaitGroup) {
			r.AddCounter(name)
			wg.Done()
		}(&wg)
	}

	wg.Wait()

	assert.Equal(t, count, len(r.Counters))

	assert.Equal(t, expectedStats, r.Stats())
}
