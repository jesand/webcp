package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/jesand/webcp/crawl"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	SW         = "webcp"
	VERSION    = "0.1.0"
	SW_VERSION = SW + " (v. " + VERSION + ")"
	USAGE      = SW_VERSION + ` - Smart site crawling

Usage:
  ` + SW + ` [-hw] <url> <dest> [--delay=<secs>] [--max-depth=<num>]
    [--resume=<path>] [--wayback-after=<date>] [--wayback-before=<date>]

All dates are in YYYY, YYYYMM, YYYYMMDD, or YYYYMMDDHHMMSS format.

Options:
  <url>                    The seed URL from which crawling should begin.
  <dest>                   The folder to which the crawl should be saved.
  --delay=<secs>           Time to wait between requests to a single domain [default: 5].
  -h --help                Show these usage notes.
  --max-depth=<num>        Stop at this tree depth [default: 5].
  --resume=<path>          Save ongoing status, and resume any previous crawls.
  --version                Show the version number.
  -w --wayback             Crawl from the Internet Wayback Machine instead.
  --wayback-after=<date>   Crawl pages archived on or after this date.
  --wayback-before=<date>  Crawl pages archived on or before this date.
`
)

func main() {
	crawler, err := ParseArgs(nil)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
	} else {
		crawler.Run()
	}
}

func ParseArgs(argv []string) (crawler crawl.Crawler, reterr error) {

	// Parse the command line and print usage
	args, _ := docopt.Parse(USAGE, argv, true, SW_VERSION, false)

	// Get the command line arguments
	var (
		delay, _               = args["--delay"].(string)
		delaySecs, delayErr    = strconv.ParseFloat(delay, 64)
		folder, _              = args["<dest>"].(string)
		depthStr, _            = args["--max-depth"].(string)
		depth, depthErr        = strconv.Atoi(depthStr)
		resume, _              = args["--resume"].(string)
		urlRaw, urlOk          = args["<url>"].(string)
		urlParsed, urlErr      = url.Parse(urlRaw)
		wayback, _             = args["--wayback"].(bool)
		wbAfter, _             = args["--wayback-after"].(string)
		wbAfterDate, wbAftErr  = ParseDate(wbAfter)
		wbBefore, _            = args["--wayback-before"].(string)
		wbBeforeDate, wbBefErr = ParseDate(wbBefore)
	)

	// Validate the command line arguments
	if delay != "" && (delayErr != nil || delaySecs < 0) {
		reterr = fmt.Errorf("Invalid --delay %q - %v", delay, delayErr)
		return
	}

	if depthStr != "" && (depthErr != nil || depth < 0) {
		reterr = fmt.Errorf("Invalid --max-depth %q", depthStr)
		return
	}

	if folder == "" {
		reterr = fmt.Errorf("<dest> is required")
		return
	} else if reterr = os.MkdirAll(folder, 0777); reterr != nil {
		return
	}

	if !urlOk || urlRaw == "" {
		reterr = fmt.Errorf("<url> is required")
		return
	} else if urlErr != nil {
		reterr = fmt.Errorf("Invalid URL %q - %v", urlRaw, urlErr)
		return
	} else if !urlParsed.IsAbs() {
		reterr = fmt.Errorf("Can't fetch non-absolute URL %s", urlRaw)
		return
	}

	if wayback {
		if wbBefore != "" && wbBefErr != nil {
			reterr = fmt.Errorf("Invalid --wayback-before date %q", wbBefore)
			return
		}
		if wbAfter != "" && wbAftErr != nil {
			reterr = fmt.Errorf("Invalid --wayback-after date %q", wbAfter)
			return
		}
		if wbBefore != "" && wbAfter != "" && !wbBeforeDate.After(wbAfterDate) {
			reterr = fmt.Errorf("--wayback-after is after --wayback-before", wbAfter)
			return
		}
	} else {
		wbAfterDate = time.Time{}
		wbBeforeDate = time.Time{}
	}

	// Build and run the crawler
	crawler = crawl.Crawler{
		FetchDelay:    time.Duration(float64(time.Second) * delaySecs),
		Folder:        folder,
		MaxDepth:      depth,
		Resume:        resume,
		Seed:          urlParsed,
		UseWayback:    wayback,
		WaybackAfter:  wbAfterDate,
		WaybackBefore: wbBeforeDate,
	}
	return
}
