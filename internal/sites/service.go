package sites

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"sort"
)

type Service interface {
	Warmup() error

	GetAll() []*Site
	GetAvailable() []*Site
	GetSortedByLatency() []*Site

	Close()
}

func NewFileSitesService(path string) Service {
	return &fileSites{filePath: path}
}

type fileSites struct {
	sites []*Site

	filePath string
}

func (s *fileSites) Warmup() error {
	file, err := s.readFile(s.filePath)
	if err != nil {
		return err
	}
	s.parseSites(file)

	return nil
}

func (s *fileSites) GetAll() []*Site {
	return s.sites
}

func (s *fileSites) GetAvailable() []*Site {
	availableSites := make([]*Site, 0, len(s.sites))

	for _, site := range s.GetAll() {
		if site.Alive {
			availableSites = append(availableSites, site)
		}
	}

	return availableSites
}

func (s *fileSites) GetSortedByLatency() []*Site {
	sites := s.GetAvailable()

	sort.Slice(sites, func(i, j int) bool {
		return sites[i].Latency < sites[j].Latency
	})

	return sites
}

// nothing to finalize
func (s *fileSites) Close() {}

func (s *fileSites) addSite(rawurl string) error {
	url, err := url.Parse(rawurl)
	if err != nil {
		return fmt.Errorf("Falied to parse %s site: %v", rawurl, err)
	}

	if url.String() == "" {
		return nil
	}

	if url.Scheme == "" {
		url.Scheme = "http"
	}

	site := &Site{
		Name: rawurl,
		Url:  url,
	}
	s.sites = append(s.sites, site)

	return nil
}

func (s *fileSites) parseSites(sites io.Reader) error {
	scanner := bufio.NewScanner(sites)
	for scanner.Scan() {
		name := scanner.Text()
		err := s.addSite(name)
		if err != nil {
			log.Printf("[ERROR] failed to parse %s site: %+v", name, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Falied to read %s sites file: %v", s.filePath, err)
	}

	return nil
}

func (s *fileSites) readFile(path string) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)

	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("Falied to open sites file: %v", err)
	}
	defer file.Close()

	io.Copy(buf, file)

	return buf, nil
}
