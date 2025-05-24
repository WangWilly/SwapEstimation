package attendance

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/dtos"
	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"github.com/WangWilly/labs-hr-go/pkgs/utils"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestGet(t *testing.T) {
	testInit(t, func(s *testSuite) {
		Convey("Given an employee with attendance records", t, func() {
			// Setup test data
			employeeID := int64(123)
			nowTime := time.Date(2023, 6, 15, 12, 0, 0, 0, time.UTC)
			attendanceID := int64(456)
			positionID := int64(789)

			// Create a mock attendance record
			attendance := &models.EmployeeAttendance{
				ID:         attendanceID,
				EmployeeID: employeeID,
				PositionID: positionID,
				ClockIn:    nowTime.Add(-4 * time.Hour), // Clocked in 4 hours ago
				ClockOut:   nowTime,                     // Clocked out now
			}

			// Expected response data structure
			expectedResponse := dtos.AttendanceV1Response{
				AttendanceID: attendance.ID,
				PositionID:   attendance.PositionID,
				ClockInTime:  utils.FormatedTime(attendance.ClockIn),
				ClockOutTime: utils.FormatedTime(attendance.ClockOut),
			}

			Convey("When retrieving the attendance by employee ID and cache hits", func() {
				// Set up cache hit expectation
				cachedResponse := &dtos.AttendanceV1Response{
					AttendanceID: attendance.ID,
					PositionID:   attendance.PositionID,
					ClockInTime:  utils.FormatedTime(attendance.ClockIn),
					ClockOutTime: utils.FormatedTime(attendance.ClockOut),
				}

				s.cacheManager.EXPECT().
					GetAttendanceV1(gomock.Any(), employeeID).
					Return(cachedResponse, nil)

				// No repository calls expected when cache hits

				// Make the request and verify response
				var actualResponse dtos.AttendanceV1Response
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/attendance/123",
					nil,
					&actualResponse,
					http.StatusOK,
				)

				Convey("Then the response should contain correct cached attendance information", func() {
					So(actualResponse.AttendanceID, ShouldEqual, expectedResponse.AttendanceID)
					So(actualResponse.PositionID, ShouldEqual, expectedResponse.PositionID)
					So(actualResponse.ClockInTime, ShouldEqual, expectedResponse.ClockInTime)
					So(actualResponse.ClockOutTime, ShouldEqual, expectedResponse.ClockOutTime)
				})
			})

			Convey("When retrieving the attendance by employee ID and cache misses", func(c C) {
				// Set up cache miss expectation
				s.cacheManager.EXPECT().
					GetAttendanceV1(gomock.Any(), employeeID).
					Return(nil, errors.New("cache miss"))

				s.employeeAttendanceRepo.EXPECT().
					Last(gomock.Any(), s.db, employeeID).
					Return(attendance, nil)

				// Add expectation that the attendance is cached before returning
				s.cacheManager.EXPECT().
					SetAttendanceV1(gomock.Any(), employeeID, gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ interface{}, _ int64, resp dtos.AttendanceV1Response, _ interface{}) error {
						// Verify the cached data matches what we expect
						c.So(resp.AttendanceID, ShouldEqual, expectedResponse.AttendanceID)
						c.So(resp.PositionID, ShouldEqual, expectedResponse.PositionID)
						c.So(resp.ClockInTime, ShouldEqual, expectedResponse.ClockInTime)
						c.So(resp.ClockOutTime, ShouldEqual, expectedResponse.ClockOutTime)
						return nil
					})

				// Make the request and verify response
				var actualResponse dtos.AttendanceV1Response
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/attendance/123",
					nil,
					&actualResponse,
					http.StatusOK,
				)

				Convey("Then the response should contain correct attendance information from database", func() {
					So(actualResponse.AttendanceID, ShouldEqual, expectedResponse.AttendanceID)
					So(actualResponse.PositionID, ShouldEqual, expectedResponse.PositionID)
					So(actualResponse.ClockInTime, ShouldEqual, expectedResponse.ClockInTime)
					So(actualResponse.ClockOutTime, ShouldEqual, expectedResponse.ClockOutTime)
				})
			})

			Convey("When retrieving attendance for employee who hasn't clocked out yet", func(c C) {
				// Create a mock attendance record where employee hasn't clocked out
				notClockedOutAttendance := &models.EmployeeAttendance{
					ID:         attendanceID,
					EmployeeID: employeeID,
					PositionID: positionID,
					ClockIn:    nowTime.Add(-4 * time.Hour), // Clocked in 4 hours ago
					ClockOut:   nowTime.Add(-4 * time.Hour), // Same as clock in (hasn't clocked out)
				}

				// Expected response with empty clock out time
				expectedNotClockedOutResponse := dtos.AttendanceV1Response{
					AttendanceID: notClockedOutAttendance.ID,
					PositionID:   notClockedOutAttendance.PositionID,
					ClockInTime:  utils.FormatedTime(notClockedOutAttendance.ClockIn),
					ClockOutTime: "", // Empty for not clocked out
				}

				s.cacheManager.EXPECT().
					GetAttendanceV1(gomock.Any(), employeeID).
					Return(nil, errors.New("cache miss"))

				s.employeeAttendanceRepo.EXPECT().
					Last(gomock.Any(), s.db, employeeID).
					Return(notClockedOutAttendance, nil)

				// Cache expectation
				s.cacheManager.EXPECT().
					SetAttendanceV1(gomock.Any(), employeeID, gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ interface{}, _ int64, resp dtos.AttendanceV1Response, _ interface{}) error {
						c.So(resp.ClockOutTime, ShouldEqual, "") // Verify empty clock out time
						return nil
					})

				// Make the request and verify response
				var actualResponse dtos.AttendanceV1Response
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/attendance/123",
					nil,
					&actualResponse,
					http.StatusOK,
				)

				Convey("Then the response should show empty clock out time", func() {
					So(actualResponse.ClockOutTime, ShouldEqual, "")
					So(actualResponse.ClockInTime, ShouldEqual, expectedNotClockedOutResponse.ClockInTime)
				})
			})

			Convey("When attendance is not found", func() {
				// Set up cache miss expectation
				s.cacheManager.EXPECT().
					GetAttendanceV1(gomock.Any(), employeeID).
					Return(nil, errors.New("cache miss"))

				// Set up expectations for failure case
				s.employeeAttendanceRepo.EXPECT().
					Last(gomock.Any(), s.db, employeeID).
					Return(nil, nil)

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/attendance/123",
					nil,
					&errorResponse,
					http.StatusNotFound,
				)

				Convey("Then the response should indicate attendance not found", func() {
					So(errorResponse["error"], ShouldEqual, "attendance not found")
				})
			})

			Convey("When there is an error retrieving attendance", func() {
				// Set up cache miss expectation
				s.cacheManager.EXPECT().
					GetAttendanceV1(gomock.Any(), employeeID).
					Return(nil, errors.New("cache miss"))

				// Set up expectations for database error
				s.employeeAttendanceRepo.EXPECT().
					Last(gomock.Any(), s.db, employeeID).
					Return(nil, errors.New("database error"))

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/attendance/123",
					nil,
					&errorResponse,
					http.StatusInternalServerError,
				)

				Convey("Then the response should indicate internal server error", func() {
					So(errorResponse["error"], ShouldEqual, "failed to get last attendance")
				})
			})
		})

		Convey("Given an invalid employee ID format", t, func() {
			Convey("When retrieving attendance with non-numeric employee ID", func() {
				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/attendance/invalid-id",
					nil,
					&errorResponse,
					http.StatusBadRequest,
				)

				Convey("Then the response should indicate invalid employee ID", func() {
					So(errorResponse["error"], ShouldEqual, "invalid employee id")
				})
			})
		})
	})
}
