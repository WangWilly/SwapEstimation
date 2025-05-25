package utils

import (
	"errors"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRetry(t *testing.T) {
	Convey("Given a retry function", t, func() {
		Convey("When the function succeeds on the first attempt", func() {
			attempts := 3
			callCount := 0

			err := retry(attempts, 1*time.Millisecond, func(i int) error {
				callCount++
				return nil
			})

			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the function should be called exactly once", func() {
				So(callCount, ShouldEqual, 1)
			})
		})

		Convey("When the function succeeds after initial failures", func() {
			attempts := 3
			callCount := 0

			err := retry(attempts, 1*time.Millisecond, func(i int) error {
				callCount++
				if callCount < 2 {
					return errors.New("temporary error")
				}
				return nil
			})

			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the function should be called exactly twice", func() {
				So(callCount, ShouldEqual, 2)
			})
		})

		Convey("When the function fails on all attempts", func() {
			attempts := 3
			callCount := 0
			expectedErr := errors.New("persistent error")

			err := retry(attempts, 1*time.Millisecond, func(i int) error {
				callCount++
				return expectedErr
			})

			Convey("Then the last error should be returned", func() {
				So(err, ShouldEqual, expectedErr)
			})

			Convey("Then the function should be called exactly three times", func() {
				So(callCount, ShouldEqual, 3)
			})
		})

		Convey("When checking the remaining attempts count passed to function", func() {
			attempts := 3
			receivedAttempts := []int{}

			_ = retry(attempts, 1*time.Millisecond, func(i int) error {
				receivedAttempts = append(receivedAttempts, i)
				return errors.New("error")
			})

			Convey("Then the attempts should decrease with each call", func() {
				So(receivedAttempts, ShouldResemble, []int{3, 2, 1})
			})
		})
	})
}
