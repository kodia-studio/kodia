package search

import (
	"errors"
	"sync"
)

/**
 * SearchManager orchestrates the search drivers.
 * Used by the system to perform searches and initiate indexing.
 */
type SearchManager struct {
	drivers map[string]SearchDriver
	defaultDriver string
	mu sync.RWMutex
}

func NewSearchManager(defaultDriver string) *SearchManager {
	return &SearchManager{
		drivers: make(map[string]SearchDriver),
		defaultDriver: defaultDriver,
	}
}

func (m *SearchManager) RegisterDriver(name string, driver SearchDriver) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.drivers[name] = driver
}

func (m *SearchManager) Driver(name ...string) (SearchDriver, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	driverName := m.defaultDriver
	if len(name) > 0 {
		driverName = name[0]
	}

	driver, ok := m.drivers[driverName]
	if !ok {
		return nil, errors.New("search driver not found: " + driverName)
	}

	return driver, nil
}

// Helper methods on manager for easier access

func (m *SearchManager) Search(index string, query string) (interface{}, error) {
	d, err := m.Driver()
	if err != nil {
		return nil, err
	}
	return d.Search(index, query, nil)
}

func (m *SearchManager) Index(index string, id string, data interface{}) error {
	d, err := m.Driver()
	if err != nil {
		return err
	}
	return d.Index(index, id, data)
}

func (m *SearchManager) Delete(index string, id string) error {
	d, err := m.Driver()
	if err != nil {
		return err
	}
	return d.Delete(index, id)
}
