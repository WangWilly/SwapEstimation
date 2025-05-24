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

func TestUpdate(t *testing.T) {
	testInit(t, func(s *testSuite) {
		Convey("Given an employee exists in the system", t, func() {
			// Setup test data
			employeeID := int64(123)
			nowTime := time.Date(2023, 6, 15, 12, 0, 0, 0, time.UTC)

			// Original employee data
			existingEmployeeInfo := &models.EmployeeInfo{
				ID:        employeeID,
				Name:      "Jane Smith",
				Age:       28,
				Address:   "456 Oak Avenue",
				Phone:     "555-5678",
				Email:     "jane.smith@example.com",
				CreatedAt: nowTime.Add(-24 * time.Hour),
				UpdatedAt: nowTime.Add(-12 * time.Hour),
			}

			// Updated employee data
			updatedInfo := UpdateRequest{
				Name:    "Jane Doe",
				Age:     29,
				Address: "789 Pine Street",
				Phone:   "555-9876",
				Email:   "jane.doe@example.com",
			}

			Convey("When updating the employee information and cache exists", func(c C) {
				// Set up expectations for successful update
				s.employeeInfoRepo.EXPECT().
					MustGet(gomock.Any(), s.db, employeeID).
					Return(existingEmployeeInfo, nil)

				s.employeeInfoRepo.EXPECT().
					Save(gomock.Any(), s.db, gomock.Any()).
					DoAndReturn(func(_ interface{}, _ interface{}, info *models.EmployeeInfo) error {
						// Check that the employee info was updated correctly
						c.So(info.ID, ShouldEqual, employeeID)
						c.So(info.Name, ShouldEqual, updatedInfo.Name)
						c.So(info.Age, ShouldEqual, updatedInfo.Age)
						c.So(info.Address, ShouldEqual, updatedInfo.Address)
						c.So(info.Phone, ShouldEqual, updatedInfo.Phone)
						c.So(info.Email, ShouldEqual, updatedInfo.Email)
						return nil
					})

					// Setup cache hit with existing employee detail
				cachedEmployeeDetail := &dtos.EmployeeV1Response{
					EmployeeID: employeeID,
					Name:       existingEmployeeInfo.Name,
					Age:        existingEmployeeInfo.Age,
					Phone:      existingEmployeeInfo.Phone,
					Email:      existingEmployeeInfo.Email,
					Address:    existingEmployeeInfo.Address,
					CreatedAt:  utils.FormatedTime(existingEmployeeInfo.CreatedAt),
					UpdatedAt:  utils.FormatedTime(existingEmployeeInfo.UpdatedAt),
					PositionID: 456,
					Position:   "Senior Developer",
					Department: "Engineering",
					Salary:     95000.00,
					StartDate:  "2023-01-01T00:00:00Z",
				}

				// Expect cache manager to be called to get existing employee details
				s.cacheManager.EXPECT().
					GetEmployeeDetailV1(gomock.Any(), employeeID).
					Return(cachedEmployeeDetail, nil)

				// Expect cache to be updated with the new employee details
				s.cacheManager.EXPECT().
					SetEmployeeDetailV1(gomock.Any(), employeeID, gomock.Any(), time.Duration(0)).
					DoAndReturn(func(_ interface{}, _ interface{}, updatedCache dtos.EmployeeV1Response, _ time.Duration) error {
						// Verify the cache was updated correctly
						c.So(updatedCache.Name, ShouldEqual, updatedInfo.Name)
						c.So(updatedCache.Age, ShouldEqual, updatedInfo.Age)
						c.So(updatedCache.Address, ShouldEqual, updatedInfo.Address)
						c.So(updatedCache.Phone, ShouldEqual, updatedInfo.Phone)
						c.So(updatedCache.Email, ShouldEqual, updatedInfo.Email)
						// Verify the other fields remain unchanged
						c.So(updatedCache.EmployeeID, ShouldEqual, cachedEmployeeDetail.EmployeeID)
						c.So(updatedCache.PositionID, ShouldEqual, cachedEmployeeDetail.PositionID)
						c.So(updatedCache.Position, ShouldEqual, cachedEmployeeDetail.Position)
						c.So(updatedCache.Department, ShouldEqual, cachedEmployeeDetail.Department)
						c.So(updatedCache.Salary, ShouldEqual, cachedEmployeeDetail.Salary)
						c.So(updatedCache.StartDate, ShouldEqual, cachedEmployeeDetail.StartDate)
						return nil
					})

				// Expected response
				expectedResponse := UpdateResponse{
					ID: employeeID,
				}

				// Make the request and verify response
				var actualResponse UpdateResponse
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPut,
					"/employee/123",
					updatedInfo,
					&actualResponse,
					http.StatusOK,
				)

				Convey("Then the response should contain the correct employee ID", func() {
					So(actualResponse.ID, ShouldEqual, expectedResponse.ID)
				})
			})

			Convey("When updating the employee information and cache misses", func(c C) {
				// Set up expectations for successful update
				s.employeeInfoRepo.EXPECT().
					MustGet(gomock.Any(), s.db, employeeID).
					Return(existingEmployeeInfo, nil)

				s.employeeInfoRepo.EXPECT().
					Save(gomock.Any(), s.db, gomock.Any()).
					DoAndReturn(func(_ interface{}, _ interface{}, info *models.EmployeeInfo) error {
						// Check that the employee info was updated correctly
						c.So(info.ID, ShouldEqual, employeeID)
						c.So(info.Name, ShouldEqual, updatedInfo.Name)
						c.So(info.Age, ShouldEqual, updatedInfo.Age)
						c.So(info.Address, ShouldEqual, updatedInfo.Address)
						c.So(info.Phone, ShouldEqual, updatedInfo.Phone)
						c.So(info.Email, ShouldEqual, updatedInfo.Email)
						return nil
					})

				// Expect cache manager to be called but return a miss
				s.cacheManager.EXPECT().
					GetEmployeeDetailV1(gomock.Any(), employeeID).
					Return(nil, errors.New("cache miss"))

				// Expected response
				expectedResponse := UpdateResponse{
					ID: employeeID,
				}

				// Make the request and verify response
				var actualResponse UpdateResponse
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPut,
					"/employee/123",
					updatedInfo,
					&actualResponse,
					http.StatusOK,
				)

				Convey("Then the response should contain the correct employee ID", func() {
					So(actualResponse.ID, ShouldEqual, expectedResponse.ID)
				})
			})

			Convey("When updating the employee information and cache update fails", func(c C) {
				// Set up expectations for successful update
				s.employeeInfoRepo.EXPECT().
					MustGet(gomock.Any(), s.db, employeeID).
					Return(existingEmployeeInfo, nil)

				s.employeeInfoRepo.EXPECT().
					Save(gomock.Any(), s.db, gomock.Any()).
					DoAndReturn(func(_ interface{}, _ interface{}, info *models.EmployeeInfo) error {
						// Check that the employee info was updated correctly
						c.So(info.ID, ShouldEqual, employeeID)
						c.So(info.Name, ShouldEqual, updatedInfo.Name)
						c.So(info.Age, ShouldEqual, updatedInfo.Age)
						c.So(info.Address, ShouldEqual, updatedInfo.Address)
						c.So(info.Phone, ShouldEqual, updatedInfo.Phone)
						c.So(info.Email, ShouldEqual, updatedInfo.Email)
						return nil
					})

				// Setup cache hit with existing employee detail
				cachedEmployeeDetail := &dtos.EmployeeV1Response{
					EmployeeID: employeeID,
					Name:       existingEmployeeInfo.Name,
					Age:        existingEmployeeInfo.Age,
					Phone:      existingEmployeeInfo.Phone,
					Email:      existingEmployeeInfo.Email,
					Address:    existingEmployeeInfo.Address,
				}

				// Expect cache manager to be called but return an error when setting
				s.cacheManager.EXPECT().
					GetEmployeeDetailV1(gomock.Any(), employeeID).
					Return(cachedEmployeeDetail, nil)

				s.cacheManager.EXPECT().
					SetEmployeeDetailV1(gomock.Any(), employeeID, gomock.Any(), time.Duration(0)).
					Return(errors.New("cache error"))

				// Expected response - should still succeed despite cache error
				expectedResponse := UpdateResponse{
					ID: employeeID,
				}

				// Make the request and verify response
				var actualResponse UpdateResponse
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPut,
					"/employee/123",
					updatedInfo,
					&actualResponse,
					http.StatusOK,
				)

				Convey("Then the response should contain the correct employee ID", func() {
					So(actualResponse.ID, ShouldEqual, expectedResponse.ID)
				})
			})

			Convey("When employee info is not found", func() {
				// Set up expectations for failure case
				s.employeeInfoRepo.EXPECT().
					MustGet(gomock.Any(), s.db, employeeID).
					Return(nil, errors.New("employee not found"))

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPut,
					"/employee/123",
					updatedInfo,
					&errorResponse,
					http.StatusNotFound,
				)

				Convey("Then the response should indicate employee not found", func() {
					So(errorResponse["error"], ShouldEqual, "employee not found")
				})
			})

			Convey("When saving employee info fails", func() {
				// Set up expectations for failure during save
				s.employeeInfoRepo.EXPECT().
					MustGet(gomock.Any(), s.db, employeeID).
					Return(existingEmployeeInfo, nil)

				s.employeeInfoRepo.EXPECT().
					Save(gomock.Any(), s.db, gomock.Any()).
					Return(errors.New("database error"))

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPut,
					"/employee/123",
					updatedInfo,
					&errorResponse,
					http.StatusInternalServerError,
				)

				Convey("Then the response should indicate a server error", func() {
					So(errorResponse["error"], ShouldEqual, "failed to update employee info")
				})
			})

			Convey("When invalid employee ID is provided", func() {
				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPut,
					"/employee/invalid",
					updatedInfo,
					&errorResponse,
					http.StatusBadRequest,
				)

				Convey("Then the response should indicate invalid ID", func() {
					So(errorResponse["error"], ShouldEqual, "invalid id")
				})
			})

			Convey("When request data is invalid", func() {
				// Invalid request with unexpected types
				invalidRequest := map[string]interface{}{
					"age": "not-a-number", // Age should be an integer
				}

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPut,
					"/employee/123",
					invalidRequest,
					&errorResponse,
					http.StatusBadRequest,
				)

				Convey("Then the response should indicate a validation error", func() {
					So(errorResponse["error"], ShouldNotBeEmpty)
				})
			})
		})
	})
}
