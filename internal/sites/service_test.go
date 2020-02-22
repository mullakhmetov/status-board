package sites

import (
	"io/ioutil"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFilesSites_WarmUp(t *testing.T) {
	path, teardown := prepFile(t)
	defer teardown()

	s := NewFileSitesService(path)

	err := s.Warmup()
	assert.NoError(t, err)

	assert.Equal(t, 4, len(s.GetAll()))
}

func TestFileSites_GetAvailable(t *testing.T) {
	path, teardown := prepFile(t)
	defer teardown()

	s := NewFileSitesService(path)

	err := s.Warmup()
	assert.NoError(t, err)

	sites := s.GetAll()
	sites[0].MarkAvailable(1)

	assert.Equal(t, 1, len(s.GetAvailable()))
}

func TestFileSites_GetSortedByLatency(t *testing.T) {
	path, teardown := prepFile(t)
	defer teardown()

	s := NewFileSitesService(path)

	err := s.Warmup()
	assert.NoError(t, err)

	sites := s.GetAll()
	for i, site := range sites {
		site.MarkAvailable(time.Duration(i))
	}

	sites = s.GetSortedByLatency()
	sorted := sort.SliceIsSorted(sites, func(i, j int) bool {
		return sites[i].Latency < sites[j].Latency
	})

	assert.True(t, sorted)
}

func prepFile(t *testing.T) (string, func()) {
	path := "/tmp/test_sites.txt"
	sites := []byte("google.com\nhttp://youtube.com\nhttps://www.facebook.com\n\ninvalid.site\n")
	err := ioutil.WriteFile(path, sites, 0644)
	assert.NoError(t, err)

	return path, func() {
		os.Remove(path)
	}
}
