package crawl

import (
	"io"
	"net/url"
)

type CrawlQueueStorage interface {
	io.Closer

	// Store a new page to crawl later
	Add(site *url.URL, depth int)

	// Get a page to crawl
	Next() (site *url.URL, depth int)
}

// Manages the crawl's frontier
type CrawlQueue struct {

	// The queue itself
	Storage CrawlQueueStorage

	// Whether we resumed a prior crawl
	DidResume bool
}

// Create a new queue
func NewQueue() CrawlQueue {
	return CrawlQueue{
		Storage: NewMemQueueStorage(),
	}
}

// Use the following resume file
func (queue *CrawlQueue) ResumeFrom(path string) error {
	var err error
	queue.Storage, queue.DidResume, err = NewFileQueueStorage(path)
	return err
}

// Store a new page to crawl later
func (queue *CrawlQueue) Add(site *url.URL, depth int) {
	loc := CanonicalURL(site)
	if !queue.Crawled(loc) {
		queue.Storage.Add(loc, depth)
	}
}

// Get a page to crawl now, blocking until one is available
func (queue *CrawlQueue) Next() (site *url.URL, depth int) {
	return queue.Storage.Next()
}

// Ask whether we've already crawled a given URL
func (queue *CrawlQueue) Crawled(site *url.URL) bool {
	// TODO: Implement Crawled()
	return false
}
