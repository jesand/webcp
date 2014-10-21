package main

import (
	"github.com/jesand/webcp/crawl"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseArgs(t *testing.T) {
	const (
		URL = "http://www.noplace.com/path/to/file.html"
	)
	var (
		URL_URL, _ = url.Parse(URL)
		tmp        string
	)

	Convey("Given just an absolute URL", t, func() {
		crawler, err := ParseArgs([]string{URL, "."})
		Convey("The correct defaults are applied", func() {
			So(err, ShouldBeNil)
			So(crawler, ShouldResemble, crawl.Crawler{
				FetchDelay:    5 * time.Second,
				Folder:        ".",
				MaxDepth:      5,
				Resume:        "",
				Seed:          URL_URL,
				UseWayback:    false,
				WaybackAfter:  time.Time{},
				WaybackBefore: time.Time{},
			})
		})
	})

	Convey("Given a relative URL", t, func() {
		_, err := ParseArgs([]string{"/path/to/file.html", "."})
		Convey("An error is returned", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given an invalid URL", t, func() {
		_, err := ParseArgs([]string{"", "."})
		Convey("An error is returned", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given a valid delay", t, func() {
		crawler, err := ParseArgs([]string{URL, ".", "--delay=1"})
		Convey("The crawler is correct", func() {
			So(err, ShouldBeNil)
			So(crawler, ShouldResemble, crawl.Crawler{
				FetchDelay:    1 * time.Second,
				Folder:        ".",
				MaxDepth:      5,
				Resume:        "",
				Seed:          URL_URL,
				UseWayback:    false,
				WaybackAfter:  time.Time{},
				WaybackBefore: time.Time{},
			})
		})
	})

	Convey("Given a non-int delay", t, func() {
		_, err := ParseArgs([]string{URL, ".", "--delay=monkey"})
		Convey("An error is returned", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given a negative delay", t, func() {
		_, err := ParseArgs([]string{URL, ".", "--delay=-1"})
		Convey("An error is returned", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given a valid new folder", t, func() {
		tmp = os.TempDir()
		Reset(func() {
			os.RemoveAll(tmp)
		})
		folder := filepath.Join(tmp, "newfolder")
		crawler, err := ParseArgs([]string{URL, folder})

		Convey("The crawler is correct", func() {
			So(err, ShouldBeNil)
			So(crawler, ShouldResemble, crawl.Crawler{
				FetchDelay:    5 * time.Second,
				Folder:        folder,
				MaxDepth:      5,
				Resume:        "",
				Seed:          URL_URL,
				UseWayback:    false,
				WaybackAfter:  time.Time{},
				WaybackBefore: time.Time{},
			})
		})

		Convey("The folder is created", func() {
			_, err := os.Stat(folder)
			So(err, ShouldBeNil)
		})
	})

	Convey("Given a valid existing folder", t, func() {
		tmp = os.TempDir()
		Reset(func() {
			os.RemoveAll(tmp)
		})
		folder := filepath.Join(tmp, "oldfolder")
		So(os.MkdirAll(folder, 0777), ShouldBeNil)
		crawler, err := ParseArgs([]string{URL, folder})

		Convey("The crawler is correct", func() {
			So(err, ShouldBeNil)
			So(crawler, ShouldResemble, crawl.Crawler{
				FetchDelay:    5 * time.Second,
				Folder:        folder,
				MaxDepth:      5,
				Resume:        "",
				Seed:          URL_URL,
				UseWayback:    false,
				WaybackAfter:  time.Time{},
				WaybackBefore: time.Time{},
			})
		})
	})

	Convey("Given an invalid folder", t, func() {
		tmp = os.TempDir()
		Reset(func() {
			os.RemoveAll(tmp)
		})
		folder := filepath.Join(tmp, "invalid:\x00folder")
		_, err := ParseArgs([]string{URL, folder})
		Convey("An error is returned", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given a positive max depth", t, func() {
		crawler, err := ParseArgs([]string{URL, ".", "--max-depth=1"})
		Convey("The crawler is correct", func() {
			So(err, ShouldBeNil)
			So(crawler, ShouldResemble, crawl.Crawler{
				FetchDelay:    5 * time.Second,
				Folder:        ".",
				MaxDepth:      1,
				Resume:        "",
				Seed:          URL_URL,
				UseWayback:    false,
				WaybackAfter:  time.Time{},
				WaybackBefore: time.Time{},
			})
		})
	})

	Convey("Given a zero max depth", t, func() {
		crawler, err := ParseArgs([]string{URL, ".", "--max-depth=0"})
		Convey("The crawler is correct", func() {
			So(err, ShouldBeNil)
			So(crawler, ShouldResemble, crawl.Crawler{
				FetchDelay:    5 * time.Second,
				Folder:        ".",
				MaxDepth:      0,
				Resume:        "",
				Seed:          URL_URL,
				UseWayback:    false,
				WaybackAfter:  time.Time{},
				WaybackBefore: time.Time{},
			})
		})
	})

	Convey("Given a non-int max depth", t, func() {
		_, err := ParseArgs([]string{URL, ".", "--max-depth=monkey"})
		Convey("An error is returned", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given a negative max depth", t, func() {
		_, err := ParseArgs([]string{URL, ".", "--max-depth=-1"})
		Convey("An error is returned", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given a non-existing resume file", t, func() {
		tmp = os.TempDir()
		Reset(func() {
			os.RemoveAll(tmp)
		})
		resume := filepath.Join(tmp, "newfile")
		crawler, err := ParseArgs([]string{URL, ".", "--resume=" + resume})
		Convey("The crawler is correct", func() {
			So(err, ShouldBeNil)
			So(crawler, ShouldResemble, crawl.Crawler{
				FetchDelay:    5 * time.Second,
				Folder:        ".",
				MaxDepth:      5,
				Resume:        resume,
				Seed:          URL_URL,
				UseWayback:    false,
				WaybackAfter:  time.Time{},
				WaybackBefore: time.Time{},
			})
		})
	})

	Convey("Given an existing resume file", t, func() {
		tmp = os.TempDir()
		Reset(func() {
			os.RemoveAll(tmp)
		})
		resume, err := ioutil.TempFile(tmp, "resume")
		So(err, ShouldBeNil)
		resume.Close()
		crawler, err := ParseArgs([]string{URL, ".", "--resume=" + resume.Name()})

		Convey("The crawler is correct", func() {
			So(err, ShouldBeNil)
			So(crawler, ShouldResemble, crawl.Crawler{
				FetchDelay:    5 * time.Second,
				Folder:        ".",
				MaxDepth:      5,
				Resume:        resume.Name(),
				Seed:          URL_URL,
				UseWayback:    false,
				WaybackAfter:  time.Time{},
				WaybackBefore: time.Time{},
			})
		})
	})

	Convey("Given --wayback", t, func() {
		crawler, err := ParseArgs([]string{URL, ".", "--wayback"})
		Convey("The crawler is correct", func() {
			So(err, ShouldBeNil)
			So(crawler, ShouldResemble, crawl.Crawler{
				FetchDelay:    5 * time.Second,
				Folder:        ".",
				MaxDepth:      5,
				Resume:        "",
				Seed:          URL_URL,
				UseWayback:    true,
				WaybackAfter:  time.Time{},
				WaybackBefore: time.Time{},
			})
		})
	})

	Convey("Given -w", t, func() {
		crawler, err := ParseArgs([]string{URL, ".", "-w"})
		Convey("The crawler is correct", func() {
			So(err, ShouldBeNil)
			So(crawler, ShouldResemble, crawl.Crawler{
				FetchDelay:    5 * time.Second,
				Folder:        ".",
				MaxDepth:      5,
				Resume:        "",
				Seed:          URL_URL,
				UseWayback:    true,
				WaybackAfter:  time.Time{},
				WaybackBefore: time.Time{},
			})
		})
	})

	Convey("Given --wayback-after and --wayback", t, func() {
		crawler, err := ParseArgs([]string{URL, ".", "--wayback", "--wayback-after=2014"})
		Convey("The crawler is correct", func() {
			So(err, ShouldBeNil)
			So(crawler, ShouldResemble, crawl.Crawler{
				FetchDelay:    5 * time.Second,
				Folder:        ".",
				MaxDepth:      5,
				Resume:        "",
				Seed:          URL_URL,
				UseWayback:    true,
				WaybackAfter:  time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC),
				WaybackBefore: time.Time{},
			})
		})
	})

	Convey("Given --wayback-after", t, func() {
		crawler, err := ParseArgs([]string{URL, ".", "--wayback-after=2014"})
		Convey("The argument is ignored", func() {
			So(err, ShouldBeNil)
			So(crawler, ShouldResemble, crawl.Crawler{
				FetchDelay:    5 * time.Second,
				Folder:        ".",
				MaxDepth:      5,
				Resume:        "",
				Seed:          URL_URL,
				UseWayback:    false,
				WaybackAfter:  time.Time{},
				WaybackBefore: time.Time{},
			})
		})
	})

	Convey("Given --wayback-before and --wayback", t, func() {
		crawler, err := ParseArgs([]string{URL, ".", "--wayback", "--wayback-before=2014"})
		Convey("The crawler is correct", func() {
			So(err, ShouldBeNil)
			So(crawler, ShouldResemble, crawl.Crawler{
				FetchDelay:    5 * time.Second,
				Folder:        ".",
				MaxDepth:      5,
				Resume:        "",
				Seed:          URL_URL,
				UseWayback:    true,
				WaybackAfter:  time.Time{},
				WaybackBefore: time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC),
			})
		})
	})

	Convey("Given --wayback-before", t, func() {
		crawler, err := ParseArgs([]string{URL, ".", "--wayback-before=2014"})
		Convey("The argument is ignored", func() {
			So(err, ShouldBeNil)
			So(crawler, ShouldResemble, crawl.Crawler{
				FetchDelay:    5 * time.Second,
				Folder:        ".",
				MaxDepth:      5,
				Resume:        "",
				Seed:          URL_URL,
				UseWayback:    false,
				WaybackAfter:  time.Time{},
				WaybackBefore: time.Time{},
			})
		})
	})

	Convey("Given --wayback-before and --wayback-after and --wayback", t, func() {
		crawler, err := ParseArgs([]string{URL, ".", "--wayback", "--wayback-after=2013", "--wayback-before=2014"})
		Convey("The crawler is correct", func() {
			So(err, ShouldBeNil)
			So(crawler, ShouldResemble, crawl.Crawler{
				FetchDelay:    5 * time.Second,
				Folder:        ".",
				MaxDepth:      5,
				Resume:        "",
				Seed:          URL_URL,
				UseWayback:    true,
				WaybackAfter:  time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC),
				WaybackBefore: time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC),
			})
		})
	})

	Convey("Given --wayback-before and --wayback-after out of order", t, func() {
		_, err := ParseArgs([]string{URL, ".", "--wayback", "--wayback-after=2014", "--wayback-before=2013"})
		Convey("An error is returned", func() {
			So(err, ShouldNotBeNil)
		})
	})
}
