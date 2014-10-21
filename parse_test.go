package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func Test_ParseDate(t *testing.T) {
	Convey("When I parse \"\"", t, func() {
		_, err := ParseDate("")
		Convey("Then I get an error", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("When I parse 14", t, func() {
		_, err := ParseDate("14")
		Convey("Then I get an error", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("When I parse 2014-02-01", t, func() {
		_, err := ParseDate("2014-02-01")
		Convey("Then I get an error", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("When I parse 2014", t, func() {
		d, err := ParseDate("2014")
		Convey("Then I get Jan 1, 2014 12:00:00 UTC", func() {
			So(err, ShouldBeNil)
			So(d, ShouldResemble, time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC))
		})
	})

	Convey("When I parse 201402", t, func() {
		d, err := ParseDate("201402")
		Convey("Then I get Feb 1, 2014 12:00:00 UTC", func() {
			So(err, ShouldBeNil)
			So(d, ShouldResemble, time.Date(2014, 2, 1, 0, 0, 0, 0, time.UTC))
		})
	})

	Convey("When I parse 20140203", t, func() {
		d, err := ParseDate("20140203")
		Convey("Then I get Feb 3, 2014 12:00:00 UTC", func() {
			So(err, ShouldBeNil)
			So(d, ShouldResemble, time.Date(2014, 2, 3, 0, 0, 0, 0, time.UTC))
		})
	})

	Convey("When I parse 20140203040506", t, func() {
		d, err := ParseDate("20140203040506")
		Convey("Then I get Feb 3, 2014 04:05:06 UTC", func() {
			So(err, ShouldBeNil)
			So(d, ShouldResemble, time.Date(2014, 2, 3, 4, 5, 6, 0, time.UTC))
		})
	})
}
