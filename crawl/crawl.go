package crawl

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// A single crawling session
type Crawler struct {

	// The root URL for the crawl
	Seed *url.URL

	// The folder to which crawled files should be stored
	Folder string

	// The maximum recursion depth
	MaxDepth int

	// The file at which to load/save session resume info
	Resume string

	// Whether to crawl using the Internet Wayback Machine
	UseWayback bool

	// Crawl the latest page that occurs within this date range
	WaybackBefore, WaybackAfter time.Time

	// The delay between successive requests to the same URL
	FetchDelay time.Duration

	// The crawler's queue
	queue CrawlQueue

	// The most recent time we accessed each recent domain
	recentDomains map[string]time.Time
}

// Run the crawl
func (crawler *Crawler) Run() {
	defer crawler.cleanup()
	if err := crawler.init(); err != nil {
		fmt.Println("Could not initialize the crawl - ", err)
		return
	}
	crawler.crawl()
}

// Initialize a new crawl
func (crawler *Crawler) init() error {

	// Set up the queue
	crawler.queue = NewQueue()
	if crawler.Resume != "" {
		if err := crawler.queue.ResumeFrom(crawler.Resume); err != nil {
			return err
		}
	}

	return nil
}

// Clean up after a crawl
func (crawler *Crawler) cleanup() {
	if crawler.queue.Storage != nil {
		crawler.queue.Storage.Close()
	}
}

// Run a crawl
func (crawler *Crawler) crawl() {

	// If we're not resuming a prior crawl, start with the seed
	if !crawler.queue.DidResume {
		crawler.queue.Add(crawler.Seed, 1)
	}

	// Crawl the frontier
	for {
		next, depth := crawler.queue.Next()
		if next == nil {
			break
		}
		crawler.fetch(next, depth, nil)
	}
}

// Fetch a URL
func (crawler *Crawler) fetch(next *url.URL, depth int, save io.Writer) {

	// Wait, if we need to
	if crawler.FetchDelay > 0 {
		noWait := time.Now().Add(-crawler.FetchDelay)
		if crawler.recentDomains == nil {
			crawler.recentDomains = make(map[string]time.Time)
		}
		lastTime, ok := crawler.recentDomains[next.Host]
		if ok && lastTime.After(noWait) {
			for key, val := range crawler.recentDomains {
				if val.Before(noWait) {
					delete(crawler.recentDomains, key)
				}
			}
			time.Sleep(crawler.FetchDelay - time.Now().Sub(lastTime))
		}
		crawler.recentDomains[next.Host] = time.Now()
	}

	// Fetch the URL
	if resp, err := http.Get(next.String()); err != nil {
		os.Stderr.WriteString("Could not fetch " + next.String() +
			" - " + err.Error())
	} else {
		if depth < crawler.MaxDepth {
			r := io.Reader(resp.Body)
			if save != nil {
				r = io.TeeReader(resp.Body, save)
			}
			crawler.parseLinks(next, r, depth+1)
		} else if save != nil {
			if _, err := io.Copy(save, resp.Body); err != nil {
				os.Stderr.WriteString("Could not save " + next.String() + " - " + err.Error())
			}
		}
		resp.Body.Close()
	}
}

// Parse a page, adding any new URLs it contains to the frontier
func (crawler *Crawler) parseLinks(source *url.URL, body io.Reader, depth int) {
	tok := html.NewTokenizer(body)
	parsing := true
	for parsing {
		switch tok.Next() {
		case html.ErrorToken:
			os.Stderr.WriteString("Could not parse " + source.String())
			parsing = false

		case html.StartTagToken:
			tag, more := tok.TagName()
			if string(tag) == "a" {
				var key, val []byte
				for more {
					key, val, more = tok.TagAttr()
					if string(key) == "href" {
						more = false
						if site, err := source.Parse(string(val)); err == nil {
							crawler.queue.Add(site, depth)
						}
					}
				}
			}
		}
	}
}
