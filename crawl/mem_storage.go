package crawl

import (
	"net/url"
)

// Memory-based storage for smaller crawls
type MemQueueStorage struct {
	Items []MemItem
}

// Create a new memory storage object
func NewMemQueueStorage() *MemQueueStorage {
	return &MemQueueStorage{}
}

type MemItem struct {
	URL   *url.URL
	Depth int
}

// Store a new page to crawl later
func (storage *MemQueueStorage) Add(site *url.URL, depth int) {
	storage.Items = append(storage.Items, MemItem{
		URL:   site,
		Depth: depth,
	})
}

// Get a page to crawl
func (storage *MemQueueStorage) Next() (*url.URL, int) {
	if len(storage.Items) == 0 {
		return nil, 0
	}
	item := storage.Items[0]
	storage.Items = storage.Items[1:]
	return item.URL, item.Depth
}

// Close the storage
func (storage *MemQueueStorage) Close() error {
	return nil
}
