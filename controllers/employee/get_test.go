package employee

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
		Convey("Given an employee exists in the system", t, func() {
			// Setup test data
			employeeID := int64(123)
			nowTime := time.Date(2023, 6, 15, 12, 0, 0, 0, time.UTC)

			employeeInfo := &models.EmployeeInfo{
				ID:        employeeID,
				Name:      "Jane Smith",
				Age:       28,
				Address:   "456 Oak Avenue",
				Phone:     "555-5678",
				Email:     "jane.smith@example.com",
				CreatedAt: nowTime.Add(-24 * time.Hour),
				UpdatedAt: nowTime.Add(-12 * time.Hour),
			}

			employeePosition := &models.EmployeePosition{
				ID:         456,
				EmployeeID: employeeID,
				Position:   "Senior Developer",
				Department: "Engineering",
				Salary:     95000.00,
				StartDate:  nowTime.Add(-6 * 30 * 24 * time.Hour), // ~6 months ago
			}

			// Expected response data structure
			expectedResponse := dtos.EmployeeV1Response{
				EmployeeID: employeeInfo.ID,
				Name:       employeeInfo.Name,
				Age:        employeeInfo.Age,
				Phone:      employeeInfo.Phone,
				Email:      employeeInfo.Email,
				Address:    employeeInfo.Address,
				CreatedAt:  utils.FormatedTime(employeeInfo.CreatedAt),
				UpdatedAt:  utils.FormatedTime(employeeInfo.UpdatedAt),
				PositionID: employeePosition.ID,
				Position:   employeePosition.Position,
				Department: employeePosition.Department,
				Salary:     employeePosition.Salary,
				StartDate:  utils.FormatedTime(employeePosition.StartDate),
			}

			Convey("When retrieving the employee by ID and cache hits", func() {
				// Set up cache hit expectation
				cachedResponse := &dtos.EmployeeV1Response{
					EmployeeID: employeeInfo.ID,
					Name:       employeeInfo.Name,
					Age:        employeeInfo.Age,
					Phone:      employeeInfo.Phone,
					Email:      employeeInfo.Email,
					Address:    employeeInfo.Address,
					CreatedAt:  utils.FormatedTime(employeeInfo.CreatedAt),
					UpdatedAt:  utils.FormatedTime(employeeInfo.UpdatedAt),
					PositionID: employeePosition.ID,
					Position:   employeePosition.Position,
					Department: employeePosition.Department,
					Salary:     employeePosition.Salary,
					StartDate:  utils.FormatedTime(employeePosition.StartDate),
				}

				s.cacheManager.EXPECT().
					GetEmployeeDetailV1(gomock.Any(), employeeID).
					Return(cachedResponse, nil)

				// No repository calls expected when cache hits

				// Make the request and verify response
				var actualResponse dtos.EmployeeV1Response
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/employee/123",
					nil,
					&actualResponse,
					http.StatusOK,
				)

				Convey("Then the response should contain correct cached employee information", func() {
					So(actualResponse.EmployeeID, ShouldEqual, expectedResponse.EmployeeID)
					So(actualResponse.Name, ShouldEqual, expectedResponse.Name)
					So(actualResponse.Age, ShouldEqual, expectedResponse.Age)
					So(actualResponse.Phone, ShouldEqual, expectedResponse.Phone)
					So(actualResponse.Email, ShouldEqual, expectedResponse.Email)
					So(actualResponse.Address, ShouldEqual, expectedResponse.Address)
					So(actualResponse.PositionID, ShouldEqual, expectedResponse.PositionID)
					So(actualResponse.Position, ShouldEqual, expectedResponse.Position)
					So(actualResponse.Department, ShouldEqual, expectedResponse.Department)
					So(actualResponse.Salary, ShouldEqual, expectedResponse.Salary)
					So(actualResponse.CreatedAt, ShouldEqual, expectedResponse.CreatedAt)
					So(actualResponse.UpdatedAt, ShouldEqual, expectedResponse.UpdatedAt)
					So(actualResponse.StartDate, ShouldEqual, expectedResponse.StartDate)
				})
			})

			Convey("When retrieving the employee by ID and cache misses", func(c C) {
				// Set up cache miss expectation
				s.cacheManager.EXPECT().
					GetEmployeeDetailV1(gomock.Any(), employeeID).
					Return(nil, errors.New("cache miss"))

				s.timeModule.EXPECT().Now().Return(nowTime)

				s.employeeInfoRepo.EXPECT().
					MustGet(gomock.Any(), s.db, employeeID).
					Return(employeeInfo, nil)

				s.employeePositionRepo.EXPECT().
					GetCurrentByEmployeeID(gomock.Any(), s.db, employeeID, nowTime).
					Return(employeePosition, nil)

					// Add expectation that the employee details are cached before returning
				s.cacheManager.EXPECT().
					SetEmployeeDetailV1(gomock.Any(), employeeID, gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ interface{}, _ int64, resp dtos.EmployeeV1Response, _ interface{}) error {
						// Verify the cached data matches what we expect
						c.So(resp.EmployeeID, ShouldEqual, expectedResponse.EmployeeID)
						c.So(resp.Name, ShouldEqual, expectedResponse.Name)
						c.So(resp.PositionID, ShouldEqual, expectedResponse.PositionID)
						return nil
					})

				// Make the request and verify response
				var actualResponse dtos.EmployeeV1Response
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/employee/123",
					nil,
					&actualResponse,
					http.StatusOK,
				)

				Convey("Then the response should contain correct employee information from database", func() {
					So(actualResponse.EmployeeID, ShouldEqual, expectedResponse.EmployeeID)
					So(actualResponse.Name, ShouldEqual, expectedResponse.Name)
					So(actualResponse.Age, ShouldEqual, expectedResponse.Age)
					So(actualResponse.Phone, ShouldEqual, expectedResponse.Phone)
					So(actualResponse.Email, ShouldEqual, expectedResponse.Email)
					So(actualResponse.Address, ShouldEqual, expectedResponse.Address)
					So(actualResponse.PositionID, ShouldEqual, expectedResponse.PositionID)
					So(actualResponse.Position, ShouldEqual, expectedResponse.Position)
					So(actualResponse.Department, ShouldEqual, expectedResponse.Department)
					So(actualResponse.Salary, ShouldEqual, expectedResponse.Salary)
					So(actualResponse.CreatedAt, ShouldEqual, expectedResponse.CreatedAt)
					So(actualResponse.UpdatedAt, ShouldEqual, expectedResponse.UpdatedAt)
					So(actualResponse.StartDate, ShouldEqual, expectedResponse.StartDate)
				})
			})

			Convey("When employee info is not found", func() {
				// Set up cache miss expectation
				s.cacheManager.EXPECT().
					GetEmployeeDetailV1(gomock.Any(), employeeID).
					Return(nil, errors.New("cache miss"))

				// Set up expectations for failure case
				s.employeeInfoRepo.EXPECT().
					MustGet(gomock.Any(), s.db, employeeID).
					Return(nil, errors.New("employee not found"))

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/employee/123",
					nil,
					&errorResponse,
					http.StatusNotFound,
				)

				Convey("Then the response should indicate employee not found", func() {
					So(errorResponse["error"], ShouldEqual, "employee not found")
				})
			})

			Convey("When employee position is not found", func() {
				// Set up cache miss expectation
				s.cacheManager.EXPECT().
					GetEmployeeDetailV1(gomock.Any(), employeeID).
					Return(nil, errors.New("cache miss"))

				// Set up expectations for partial failure
				s.employeeInfoRepo.EXPECT().
					MustGet(gomock.Any(), s.db, employeeID).
					Return(employeeInfo, nil)

				s.timeModule.EXPECT().Now().Return(nowTime)

				s.employeePositionRepo.EXPECT().
					GetCurrentByEmployeeID(gomock.Any(), s.db, employeeID, nowTime).
					Return(nil, errors.New("employee position not found"))

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/employee/123",
					nil,
					&errorResponse,
					http.StatusNotFound,
				)

				Convey("Then the response should indicate position not found", func() {
					So(errorResponse["error"], ShouldEqual, "employee position not found")
				})
			})
		})

		Convey("Given an invalid employee ID format", t, func() {
			Convey("When retrieving the employee with non-numeric ID", func() {
				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/employee/invalid-id",
					nil,
					&errorResponse,
					http.StatusBadRequest,
				)

				Convey("Then the response should indicate invalid ID", func() {
					So(errorResponse["error"], ShouldEqual, "invalid id")
				})
			})
		})
	})
}
