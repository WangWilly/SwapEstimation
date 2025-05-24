package utils

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFormatedTime(t *testing.T) {
	Convey("Given a FormatedTime function", t, func() {
		Convey("When formatting a UTC time", func() {
			testTime := time.Date(2023, 6, 15, 12, 30, 45, 0, time.UTC)
			formattedTime := FormatedTime(testTime)

			Convey("Then it should return the correct formatted string", func() {
				So(formattedTime, ShouldEqual, "2023-06-15 12:30:45")
			})
		})

		Convey("When formatting a non-UTC time", func() {
			// Create a time in a different timezone (UTC+8)
			loc, _ := time.LoadLocation("Asia/Taipei")
			testTime := time.Date(2023, 6, 15, 20, 30, 45, 0, loc)
			formattedTime := FormatedTime(testTime)

			Convey("Then it should convert to UTC and return the correct formatted string", func() {
				// 20:30:45 in UTC+8 is 12:30:45 in UTC
				So(formattedTime, ShouldEqual, "2023-06-15 12:30:45")
			})
		})

		Convey("When formatting a zero time value", func() {
			testTime := time.Time{}
			formattedTime := FormatedTime(testTime)

			Convey("Then it should format the zero time in UTC", func() {
				So(formattedTime, ShouldEqual, "0001-01-01 00:00:00")
			})
		})
	})
}
