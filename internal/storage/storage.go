package storage

import (
	"errors"
	"sync"
	"time"
)

var ErrSiteNotFound = errors.New("site not found")

type Site struct {
	ID           int
	URL          string
	Availability bool
	LastCheck    time.Time
	ResponseTime time.Duration
}

type Storage struct {
	Sites  map[int]Site
	nextID int
	mu     sync.RWMutex
}

func NewStorage() *Storage {
	sites := make(map[int]Site)
	return &Storage{
		Sites:  sites,
		nextID: 0,
	}
}

func (s *Storage) GetAll() []Site {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sites := make([]Site, 0, len(s.Sites))

	for _, site := range s.Sites {
		sites = append(sites, site)
	}
	return sites
}

func (s *Storage) Add(site Site) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := ValidateURL(site.URL)
	if err != nil {
		return err
	}

	s.nextID++
	site.ID = s.nextID
	s.Sites[site.ID] = site

	return nil
}

func (s *Storage) GetByID(id int) (Site, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	site, exists := s.Sites[id]
	if !exists {
		return Site{}, ErrSiteNotFound
	}
	return site, nil
}

func (s *Storage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	site, exists := s.Sites[id]
	if !exists {
		return ErrSiteNotFound
	}

	delete(s.Sites, site.ID)
	return nil
}

func (s *Storage) Update(id int, newURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	site, exists := s.Sites[id]
	if !exists {
		return ErrSiteNotFound
	}

	if err := ValidateURL(newURL); err != nil {
		return err
	}

	site.URL = newURL
	s.Sites[id] = site
	return nil
}
