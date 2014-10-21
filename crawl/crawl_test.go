package crawl

import (
	"bytes"
	"code.google.com/p/gomock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

const (
	NO_LINK_PAGE = `<html><body>Liberating, isn't it?</body></html>`

	REL_LINK_PAGE = `<html><body>
<a href="/some/page.html">anchor</a>
<a href="/some/../page3.html">anchor</a>
<a href="page2.html">anchor</a>
</body></html>`

	ABS_LINK_PAGE = `<html><body>
<a href="http://domain.com/some/page.html">anchor</a>
<a href="http://domain2.com">anchor</a>
</body></html>`
)

var (
	REL_LINKS_STR = []string{
		"/some/page.html",
		"/some/../page3.html",
		"page2.html",
	}

	ABS_LINKS_STR = []string{
		"http://domain.com/some/page.html",
		"http://domain2.com",
	}
)

type Handler struct {
	Next string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	buff := bytes.NewBufferString(h.Next)
	io.Copy(w, buff)
}

func TestCrawler(t *testing.T) {
	var REL_LINKS, ABS_LINKS []*url.URL
	for _, s := range ABS_LINKS_STR {
		u, _ := url.Parse(s)
		ABS_LINKS = append(ABS_LINKS, u)
	}

	Convey("Given a crawler with mock queue storage", t, func() {
		var handler Handler
		srv := httptest.NewServer(&handler)
		So(srv, ShouldNotBeNil)
		srvURL, err := url.Parse(srv.URL)
		So(err, ShouldBeNil)

		REL_LINKS = nil
		for _, s := range REL_LINKS_STR {
			u, _ := url.Parse(s)
			REL_LINKS = append(REL_LINKS, srvURL.ResolveReference(u))
		}

		Reset(func() {
			srv.Close()
		})

		ctrl := gomock.NewController(t)
		storage := NewMockCrawlQueueStorage(ctrl)
		Reset(func() {
			ctrl.Finish()
		})
		crawler := Crawler{
			MaxDepth: 5,
			queue:    NewQueue(),
		}
		crawler.queue.Storage = storage

		Convey("When I fetch a page I want to save", func() {

			// Then I add the links
			gomock.InOrder(
				storage.EXPECT().Add(ABS_LINKS[0], 2),
				storage.EXPECT().Add(ABS_LINKS[1], 2),
			)

			buff := bytes.Buffer{}
			handler.Next = ABS_LINK_PAGE
			crawler.fetch(srvURL, 1, &buff)

			Convey("Then I save the page", func() {
				So(string(buff.Bytes()), ShouldEqual, ABS_LINK_PAGE)
			})
		})

		Convey("When I fetch a page I don't want to save", func() {

			// Then I add the links
			gomock.InOrder(
				storage.EXPECT().Add(ABS_LINKS[0], 2),
				storage.EXPECT().Add(ABS_LINKS[1], 2),
			)

			handler.Next = ABS_LINK_PAGE
			crawler.fetch(srvURL, 1, nil)
		})

		Convey("When I fetch a page at the maximum depth", func() {

			// Then I don't add the links
			buff := bytes.Buffer{}
			handler.Next = ABS_LINK_PAGE
			crawler.fetch(srvURL, crawler.MaxDepth, &buff)

			Convey("Then I save the page", func() {
				So(string(buff.Bytes()), ShouldEqual, ABS_LINK_PAGE)
			})
		})

		Convey("When I parse a page with no links", func() {

			// Then I don't add the links
			buff := bytes.Buffer{}
			handler.Next = NO_LINK_PAGE
			crawler.fetch(srvURL, 1, &buff)

			Convey("Then I save the page", func() {
				So(string(buff.Bytes()), ShouldEqual, NO_LINK_PAGE)
			})
		})

		Convey("When I parse a page with relative links", func() {

			// Then I add the links
			gomock.InOrder(
				storage.EXPECT().Add(REL_LINKS[0], 2),
				storage.EXPECT().Add(REL_LINKS[1], 2),
				storage.EXPECT().Add(REL_LINKS[2], 2),
			)

			handler.Next = REL_LINK_PAGE
			crawler.fetch(srvURL, 1, nil)
		})

		Convey("When I fetch a page I need to wait for", func() {
			handler.Next = NO_LINK_PAGE
			crawler.FetchDelay = time.Millisecond * 250
			crawler.recentDomains = make(map[string]time.Time)
			crawler.recentDomains[srvURL.Host] = time.Now().Add(-time.Millisecond * 100)
			before := time.Now()
			crawler.fetch(srvURL, 1, nil)
			after := time.Now()

			Convey("Then I wait for the delay period", func() {
				So(after.Sub(before), ShouldBeBetween, 150*time.Millisecond, 160*time.Millisecond)
			})
		})

		Convey("When I fetch a page I don't need to wait for", func() {
			handler.Next = NO_LINK_PAGE
			crawler.FetchDelay = time.Millisecond * 250
			before := time.Now()
			crawler.fetch(srvURL, 1, nil)
			after := time.Now()

			Convey("Then I don't wait for the delay period", func() {
				So(after.Sub(before), ShouldBeLessThan, time.Millisecond)
			})
		})
	})
}
