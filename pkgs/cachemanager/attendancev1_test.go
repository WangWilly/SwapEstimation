package cachemanager

import (
	"testing"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/dtos"
	. "github.com/smartystreets/goconvey/convey"
)

////////////////////////////////////////////////////////////////////////////////

func TestAttendanceV1(t *testing.T) {
	testInit(t, func(s *testSuite) {
		Convey("Given an attendance cache manager", t, func() {
			ctx := t.Context()
			employeeID := int64(123)

			// Create test data
			attendanceData := dtos.AttendanceV1Response{
				AttendanceID: int64(789),
				PositionID:   int64(456),
				ClockInTime:  "2023-06-15 09:00:00",
				ClockOutTime: "2023-06-15 17:00:00",
			}

			// Clean up before testing to ensure consistent state
			_ = s.manager.DeleteAttendanceV1(ctx, employeeID)

			Convey("When setting attendance cache", func() {
				err := s.manager.SetAttendanceV1(ctx, employeeID, attendanceData, time.Minute*15)

				Convey("Then no error should occur", func() {
					So(err, ShouldBeNil)
				})

				Convey("And when getting the cached attendance", func() {
					cachedData, err := s.manager.GetAttendanceV1(ctx, employeeID)

					Convey("Then no error should occur", func() {
						So(err, ShouldBeNil)
						So(cachedData, ShouldNotBeNil)
					})

					Convey("And the cached data should match the original data", func() {
						So(cachedData.AttendanceID, ShouldEqual, attendanceData.AttendanceID)
						So(cachedData.PositionID, ShouldEqual, attendanceData.PositionID)
						So(cachedData.ClockInTime, ShouldEqual, attendanceData.ClockInTime)
						So(cachedData.ClockOutTime, ShouldEqual, attendanceData.ClockOutTime)
					})
				})
			})

			Convey("When updating existing attendance cache", func() {
				// First set the initial data
				err := s.manager.SetAttendanceV1(ctx, employeeID, attendanceData, time.Minute*15)
				So(err, ShouldBeNil)

				// Update the data
				updatedData := attendanceData
				updatedData.ClockOutTime = "2023-06-15 18:30:00"

				// Set the updated data
				err = s.manager.SetAttendanceV1(ctx, employeeID, updatedData, time.Minute*15)

				Convey("Then no error should occur", func() {
					So(err, ShouldBeNil)
				})

				Convey("And when getting the updated cached attendance", func() {
					cachedData, err := s.manager.GetAttendanceV1(ctx, employeeID)

					Convey("Then the cached data should contain the updates", func() {
						So(err, ShouldBeNil)
						So(cachedData, ShouldNotBeNil)
						So(cachedData.AttendanceID, ShouldEqual, updatedData.AttendanceID)
						So(cachedData.PositionID, ShouldEqual, updatedData.PositionID)
						So(cachedData.ClockInTime, ShouldEqual, updatedData.ClockInTime)
						So(cachedData.ClockOutTime, ShouldEqual, updatedData.ClockOutTime)
					})
				})
			})

			Convey("When getting a non-existent attendance", func() {
				nonExistentID := int64(999)
				// Make sure it doesn't exist
				_ = s.manager.DeleteAttendanceV1(ctx, nonExistentID)

				cachedData, err := s.manager.GetAttendanceV1(ctx, nonExistentID)

				Convey("Then an error should not occur", func() {
					So(err, ShouldBeNil)
				})

				Convey("And the cached data should be nil", func() {
					So(cachedData, ShouldBeNil)
				})
			})

			Convey("When deleting an attendance cache", func() {
				// First set the data
				err := s.manager.SetAttendanceV1(ctx, employeeID, attendanceData, time.Minute*15)
				So(err, ShouldBeNil)

				// Then delete it
				err = s.manager.DeleteAttendanceV1(ctx, employeeID)

				Convey("Then an error should not occur", func() {
					So(err, ShouldBeNil)
				})

				Convey("And when trying to get the deleted cache", func() {
					cachedData, err := s.manager.GetAttendanceV1(ctx, employeeID)

					Convey("Then an error should not occur", func() {
						So(err, ShouldBeNil)
					})

					Convey("And the cached data should be nil", func() {
						So(cachedData, ShouldBeNil)
					})
				})
			})

			Convey("When deleting a non-existent attendance cache", func() {
				nonExistentID := int64(999)
				// Make sure it doesn't exist
				_ = s.manager.DeleteAttendanceV1(ctx, nonExistentID)

				// Try to delete it again
				err := s.manager.DeleteAttendanceV1(ctx, nonExistentID)

				Convey("Then no error should occur", func() {
					So(err, ShouldBeNil)
				})
			})

			Convey("When performing the full lifecycle of attendance cache operations", func() {
				// 1. Set initial data
				err := s.manager.SetAttendanceV1(ctx, employeeID, attendanceData, time.Minute*15)
				So(err, ShouldBeNil)

				// 2. Get the data to verify it was set
				cachedData, err := s.manager.GetAttendanceV1(ctx, employeeID)
				So(err, ShouldBeNil)
				So(cachedData, ShouldNotBeNil)
				So(cachedData.AttendanceID, ShouldEqual, attendanceData.AttendanceID)

				// 3. Delete the data
				err = s.manager.DeleteAttendanceV1(ctx, employeeID)
				So(err, ShouldBeNil)

				// 4. Verify that the data was deleted
				cachedData, err = s.manager.GetAttendanceV1(ctx, employeeID)
				So(err, ShouldBeNil)
				So(cachedData, ShouldBeNil)

				// 5. Set the data again with no expiration (0 duration)
				err = s.manager.SetAttendanceV1(ctx, employeeID, attendanceData, 0)
				So(err, ShouldBeNil)

				// 6. Verify it was set again
				cachedData, err = s.manager.GetAttendanceV1(ctx, employeeID)
				So(err, ShouldBeNil)
				So(cachedData, ShouldNotBeNil)
				So(cachedData.AttendanceID, ShouldEqual, attendanceData.AttendanceID)
			})
		})
	})
}
