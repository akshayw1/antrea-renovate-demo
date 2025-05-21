package storage

import (
	"sync"
)

// Item represents a stored item
type Item struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

// MemoryStorage provides a simple in-memory storage
type MemoryStorage struct {
	items map[string]Item
	mu    sync.RWMutex
}

// NewMemoryStorage creates a new memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		items: make(map[string]Item),
	}
}

// GetItem retrieves an item by ID
func (s *MemoryStorage) GetItem(id string) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	item, exists := s.items[id]
	return item, exists
}

// ListItems returns all stored items
func (s *MemoryStorage) ListItems() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	items := make([]Item, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}
	return items
}

// StoreItem adds or updates an item
func (s *MemoryStorage) StoreItem(item Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.items[item.ID] = item
}

// DeleteItem removes an item by ID
func (s *MemoryStorage) DeleteItem(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.items[id]; exists {
		delete(s.items, id)
		return true
	}
	return false
}