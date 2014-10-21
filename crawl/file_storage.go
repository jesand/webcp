package crawl

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// File-based queue storage for resumable crawls
type FileQueueStorage struct {
	Reader  *os.File
	Scanner *bufio.Scanner
	Writer  *os.File
}

// Open/Create a new file storage at a given path
func NewFileQueueStorage(path string) (storage *FileQueueStorage, didResume bool, err error) {
	storage = &FileQueueStorage{}
	defer func() {
		if err != nil {
			storage.Close()
		}
	}()

	// Open the file for writing
	if storage.Writer, err = os.OpenFile(path, os.O_APPEND, 0644); err != nil {
		return
	}

	// Find the last crawled page in the file
	var last string
	r, err := os.Open(path)
	if err != nil {
		return nil, false, err
	} else {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "- ") {
				last = line[2:]
			}
		}
		r.Close()
	}

	// Seek to the last crawled page
	if storage.Reader, err = os.Open(path); err != nil {
		return
	}
	storage.Scanner = bufio.NewScanner(r)
	if last != "" {
		didResume = true
		for storage.Scanner.Scan() {
			line := storage.Scanner.Text()
			if line == last {
				break
			}
		}
	}
	return
}

// Store a new page to crawl later
func (storage *FileQueueStorage) Add(site *url.URL, depth int) {
	line := fmt.Sprintf("%d %s\n", depth, site.String())
	if _, err := storage.Writer.WriteString(line); err != nil {
		os.Stderr.WriteString("Failed to record " + site.String() +
			" - " + err.Error() + "\n")
	}
}

// Get a page to crawl now, blocking until one is available
func (storage *FileQueueStorage) Next() (*url.URL, int) {
	for storage.Scanner.Scan() {
		line := storage.Scanner.Text()
		if !strings.HasPrefix(line, "- ") {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) != 2 {
				os.Stderr.WriteString("Invalid line in restore file\n")
			} else if depth, err := strconv.Atoi(parts[0]); err != nil {
				os.Stderr.WriteString("Invalid depth field in restore file\n")
			} else if site, err := url.Parse(parts[1]); err != nil {
				os.Stderr.WriteString("Invalid URL: " + parts[1] + "\n")
			} else {
				_, err := storage.Writer.WriteString("- " + parts[1] + "\n")
				if err != nil {
					os.Stderr.WriteString("Failed to record crawl for " +
						parts[1] + " - " + err.Error() + "\n")
				}
				return site, depth
			}
		}
	}
	return nil, 0
}

// Close the underlying files
func (storage *FileQueueStorage) Close() error {
	if storage.Reader != nil {
		storage.Reader.Close()
		storage.Reader = nil
		storage.Scanner = nil
	}
	if storage.Writer != nil {
		storage.Writer.Close()
		storage.Writer = nil
	}
	return nil
}
