package attendance

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/dtos"
	"github.com/WangWilly/labs-hr-go/pkgs/models"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	testInit(t, func(s *testSuite) {
		Convey("Given an employee exists in the system", t, func() {
			// Setup test data
			employeeID := int64(123)
			positionID := int64(456)
			attendanceID := int64(789)
			nowTime := time.Date(2023, 6, 15, 9, 0, 0, 0, time.UTC)

			employeePosition := &models.EmployeePosition{
				ID:         positionID,
				EmployeeID: employeeID,
				Position:   "Software Engineer",
				Department: "Engineering",
				Salary:     90000.00,
				StartDate:  nowTime.Add(-30 * 24 * time.Hour), // 30 days ago
			}

			req := CreateRequest{
				EmployeeID: employeeID,
			}

			Convey("When clocking in for the first time", func() {
				// New attendance record for clock-in
				newAttendance := &models.EmployeeAttendance{
					ID:         attendanceID,
					EmployeeID: employeeID,
					PositionID: positionID,
					ClockIn:    nowTime,
					// ClockOut is zero time
				}

				// Set up expectations
				s.timeModule.EXPECT().Now().Return(nowTime).Times(2)

				// Add cache manager expectations
				s.cacheManager.EXPECT().
					GetEmployeeDetailV1(gomock.Any(), employeeID).
					Return(nil, nil)

				s.employeePositionRepo.EXPECT().
					GetCurrentByEmployeeID(gomock.Any(), s.db, employeeID, nowTime).
					Return(employeePosition, nil)

				s.employeeAttendanceRepo.EXPECT().
					Last(gomock.Any(), s.db, employeeID).
					Return(nil, nil)

				s.employeeAttendanceRepo.EXPECT().
					CreateForClockIn(gomock.Any(), s.db, employeeID, positionID, nowTime).
					Return(newAttendance, nil)

				// Expected response for cache
				expectedResponse := dtos.AttendanceV1Response{
					AttendanceID: attendanceID,
					PositionID:   positionID,
					ClockInTime:  "2023-06-15 09:00:00",
					ClockOutTime: "",
				}

				// Expect cache set call
				s.cacheManager.EXPECT().
					SetAttendanceV1(gomock.Any(), employeeID, expectedResponse, gomock.Any()).
					Return(nil)

				// Make the request and verify response
				var actualResponse dtos.AttendanceV1Response
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPost,
					"/attendance",
					req,
					&actualResponse,
					http.StatusCreated,
				)

				Convey("Then the response should contain the clock-in details", func() {
					So(actualResponse.AttendanceID, ShouldEqual, expectedResponse.AttendanceID)
					So(actualResponse.PositionID, ShouldEqual, expectedResponse.PositionID)
					So(actualResponse.ClockInTime, ShouldEqual, expectedResponse.ClockInTime)
					So(actualResponse.ClockOutTime, ShouldEqual, expectedResponse.ClockOutTime)
				})
			})

			Convey("When clocking out after previous clock-in", func() {
				clockInTime := nowTime.Add(-8 * time.Hour) // 8 hours before now
				existingAttendance := &models.EmployeeAttendance{
					ID:         attendanceID,
					EmployeeID: employeeID,
					PositionID: positionID,
					ClockIn:    clockInTime,
					ClockOut:   clockInTime, // Same as ClockIn indicating no clock-out yet
				}

				updatedAttendance := &models.EmployeeAttendance{
					ID:         attendanceID,
					EmployeeID: employeeID,
					PositionID: positionID,
					ClockIn:    clockInTime,
					ClockOut:   nowTime,
				}

				// Set up expectations
				s.timeModule.EXPECT().Now().Return(nowTime).Times(2)

				// Add cache manager expectations
				s.cacheManager.EXPECT().
					GetEmployeeDetailV1(gomock.Any(), employeeID).
					Return(nil, nil)

				s.employeePositionRepo.EXPECT().
					GetCurrentByEmployeeID(gomock.Any(), s.db, employeeID, nowTime).
					Return(employeePosition, nil)

				s.employeeAttendanceRepo.EXPECT().
					Last(gomock.Any(), s.db, employeeID).
					Return(existingAttendance, nil)

				s.employeeAttendanceRepo.EXPECT().
					UpdateForClockOut(gomock.Any(), s.db, attendanceID, nowTime).
					Return(updatedAttendance, nil)

				// Expected response for cache
				expectedResponse := dtos.AttendanceV1Response{
					AttendanceID: attendanceID,
					PositionID:   positionID,
					ClockInTime:  "2023-06-15 01:00:00",
					ClockOutTime: "2023-06-15 09:00:00",
				}

				// Expect cache set call
				s.cacheManager.EXPECT().
					SetAttendanceV1(gomock.Any(), employeeID, expectedResponse, gomock.Any()).
					Return(nil)

				// Make the request and verify response
				var actualResponse dtos.AttendanceV1Response
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPost,
					"/attendance",
					req,
					&actualResponse,
					http.StatusCreated,
				)

				Convey("Then the response should contain both clock-in and clock-out times", func() {
					So(actualResponse.AttendanceID, ShouldEqual, expectedResponse.AttendanceID)
					So(actualResponse.PositionID, ShouldEqual, expectedResponse.PositionID)
					So(actualResponse.ClockInTime, ShouldEqual, expectedResponse.ClockInTime)
					So(actualResponse.ClockOutTime, ShouldEqual, expectedResponse.ClockOutTime)
				})
			})

			Convey("When employee position is not found", func() {
				// Set up expectations for failure
				s.timeModule.EXPECT().Now().Return(nowTime)

				s.cacheManager.EXPECT().
					GetEmployeeDetailV1(gomock.Any(), employeeID).
					Return(nil, nil)

				s.employeePositionRepo.EXPECT().
					GetCurrentByEmployeeID(gomock.Any(), s.db, employeeID, nowTime).
					Return(nil, nil)

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPost,
					"/attendance",
					req,
					&errorResponse,
					http.StatusInternalServerError,
				)

				Convey("Then the response should indicate position not found", func() {
					So(errorResponse["error"], ShouldEqual, "failed to get employee position")
				})
			})

			Convey("When getting position fails", func() {
				// Set up expectations for failure
				s.cacheManager.EXPECT().
					GetEmployeeDetailV1(gomock.Any(), employeeID).
					Return(nil, errors.New("cache error"))

				s.timeModule.EXPECT().Now().Return(nowTime)

				s.employeePositionRepo.EXPECT().
					GetCurrentByEmployeeID(gomock.Any(), s.db, employeeID, nowTime).
					Return(nil, errors.New("database error"))

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPost,
					"/attendance",
					req,
					&errorResponse,
					http.StatusInternalServerError,
				)

				Convey("Then the response should indicate a server error", func() {
					So(errorResponse["error"], ShouldEqual, "failed to get employee position")
				})
			})

			Convey("When getting last attendance fails", func() {
				// Set up expectations for failure
				s.timeModule.EXPECT().Now().Return(nowTime)

				s.cacheManager.EXPECT().
					GetEmployeeDetailV1(gomock.Any(), employeeID).
					Return(nil, nil)

				s.employeePositionRepo.EXPECT().
					GetCurrentByEmployeeID(gomock.Any(), s.db, employeeID, nowTime).
					Return(employeePosition, nil)

				s.employeeAttendanceRepo.EXPECT().
					Last(gomock.Any(), s.db, employeeID).
					Return(nil, errors.New("database error"))

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPost,
					"/attendance",
					req,
					&errorResponse,
					http.StatusInternalServerError,
				)

				Convey("Then the response should indicate a server error", func() {
					So(errorResponse["error"], ShouldEqual, "failed to get employee attendance")
				})
			})

			Convey("When creating a new attendance record fails", func() {
				// Set up expectations for failure
				s.timeModule.EXPECT().Now().Return(nowTime).Times(2)

				s.cacheManager.EXPECT().
					GetEmployeeDetailV1(gomock.Any(), employeeID).
					Return(nil, nil)

				s.employeePositionRepo.EXPECT().
					GetCurrentByEmployeeID(gomock.Any(), s.db, employeeID, nowTime).
					Return(employeePosition, nil)

				s.employeeAttendanceRepo.EXPECT().
					Last(gomock.Any(), s.db, employeeID).
					Return(nil, nil)

				s.employeeAttendanceRepo.EXPECT().
					CreateForClockIn(gomock.Any(), s.db, employeeID, positionID, nowTime).
					Return(nil, errors.New("database error"))

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPost,
					"/attendance",
					req,
					&errorResponse,
					http.StatusInternalServerError,
				)

				Convey("Then the response should indicate a server error", func() {
					So(errorResponse["error"], ShouldEqual, "failed to create/update attendance")
				})
			})

			Convey("When updating attendance for clock-out fails", func() {
				clockInTime := nowTime.Add(-8 * time.Hour)
				existingAttendance := &models.EmployeeAttendance{
					ID:         attendanceID,
					EmployeeID: employeeID,
					PositionID: positionID,
					ClockIn:    clockInTime,
					ClockOut:   clockInTime,
				}

				// Set up expectations for failure during update
				s.timeModule.EXPECT().Now().Return(nowTime).Times(2)

				s.cacheManager.EXPECT().
					GetEmployeeDetailV1(gomock.Any(), employeeID).
					Return(nil, nil)

				s.employeePositionRepo.EXPECT().
					GetCurrentByEmployeeID(gomock.Any(), s.db, employeeID, nowTime).
					Return(employeePosition, nil)

				s.employeeAttendanceRepo.EXPECT().
					Last(gomock.Any(), s.db, employeeID).
					Return(existingAttendance, nil)

				s.employeeAttendanceRepo.EXPECT().
					UpdateForClockOut(gomock.Any(), s.db, attendanceID, nowTime).
					Return(nil, errors.New("database error"))

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPost,
					"/attendance",
					req,
					&errorResponse,
					http.StatusInternalServerError,
				)

				Convey("Then the response should indicate a server error", func() {
					So(errorResponse["error"], ShouldEqual, "failed to create/update attendance")
				})
			})

			Convey("When setting attendance to cache fails", func() {
				// New attendance record for clock-in
				newAttendance := &models.EmployeeAttendance{
					ID:         attendanceID,
					EmployeeID: employeeID,
					PositionID: positionID,
					ClockIn:    nowTime,
				}

				// Set up expectations
				s.timeModule.EXPECT().Now().Return(nowTime).Times(2)

				s.cacheManager.EXPECT().
					GetEmployeeDetailV1(gomock.Any(), employeeID).
					Return(nil, nil)

				s.employeePositionRepo.EXPECT().
					GetCurrentByEmployeeID(gomock.Any(), s.db, employeeID, nowTime).
					Return(employeePosition, nil)

				s.employeeAttendanceRepo.EXPECT().
					Last(gomock.Any(), s.db, employeeID).
					Return(nil, nil)

				s.employeeAttendanceRepo.EXPECT().
					CreateForClockIn(gomock.Any(), s.db, employeeID, positionID, nowTime).
					Return(newAttendance, nil)

				expectedResponse := dtos.AttendanceV1Response{
					AttendanceID: attendanceID,
					PositionID:   positionID,
					ClockInTime:  "2023-06-15 09:00:00",
					ClockOutTime: "",
				}

				// Expect cache set to fail but API should still succeed
				s.cacheManager.EXPECT().
					SetAttendanceV1(gomock.Any(), employeeID, expectedResponse, gomock.Any()).
					Return(errors.New("cache error"))

				var actualResponse dtos.AttendanceV1Response
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPost,
					"/attendance",
					req,
					&actualResponse,
					http.StatusCreated,
				)

				Convey("Then the API should still succeed even if caching fails", func() {
					So(actualResponse.AttendanceID, ShouldEqual, expectedResponse.AttendanceID)
					So(actualResponse.PositionID, ShouldEqual, expectedResponse.PositionID)
					So(actualResponse.ClockInTime, ShouldEqual, expectedResponse.ClockInTime)
					So(actualResponse.ClockOutTime, ShouldEqual, expectedResponse.ClockOutTime)
				})
			})
		})
	})
}
