package employee

import (
	"net/http"
	"testing"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/dtos"
	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"github.com/WangWilly/labs-hr-go/pkgs/utils"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	testInit(t, func(s *testSuite) {
		Convey("Given a create employee request", t, func() {
			// Setup test data
			employeeInfo := &models.EmployeeInfo{
				ID:      1,
				Name:    "John Doe",
				Age:     30,
				Address: "123 Main St",
				Phone:   "555-1234",
				Email:   "john.doe@example.com",
			}

			startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
			nowTime := time.Date(2023, 1, 15, 12, 0, 0, 0, time.UTC)

			employeePosition := &models.EmployeePosition{
				ID:         2,
				EmployeeID: employeeInfo.ID,
				Position:   "Developer",
				Department: "Engineering",
				Salary:     75000,
				StartDate:  time.Unix(startDate, 0),
			}

			// Create request payload
			req := CreateRequest{
				Name:       employeeInfo.Name,
				Age:        employeeInfo.Age,
				Address:    employeeInfo.Address,
				Phone:      employeeInfo.Phone,
				Email:      employeeInfo.Email,
				Position:   employeePosition.Position,
				Department: employeePosition.Department,
				Salary:     employeePosition.Salary,
				StartDate:  startDate,
			}

			Convey("When creating a new employee", func() {
				// Set up expectations
				s.timeModule.EXPECT().Now().Return(nowTime)

				s.employeeInfoRepo.EXPECT().
					Create(gomock.Any(), s.db, gomock.Any()).
					DoAndReturn(func(_ interface{}, _ interface{}, info *models.EmployeeInfo) error {
						info.ID = employeeInfo.ID
						return nil
					})

				s.employeePositionRepo.EXPECT().
					Create(gomock.Any(), s.db, gomock.Any(), nowTime).
					DoAndReturn(func(_ interface{}, _ interface{}, pos *models.EmployeePosition, _ time.Time) error {
						pos.ID = employeePosition.ID
						return nil
					})

					// Expect cache manager to be called with the correct employee details
				expectedCache := dtos.EmployeeV1Response{
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
					SetEmployeeDetailV1(gomock.Any(), employeeInfo.ID, gomock.Eq(expectedCache), time.Duration(0)).
					Return(nil)

				// Expected response
				expectedResponse := CreateResponse{
					EmployeeID: employeeInfo.ID,
					PositionID: employeePosition.ID,
				}

				// Make the request and verify response
				var actualResponse CreateResponse
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPost,
					"/employee",
					req,
					&actualResponse,
					http.StatusCreated,
				)

				Convey("Then the response should contain the correct IDs", func() {
					So(actualResponse.EmployeeID, ShouldEqual, expectedResponse.EmployeeID)
					So(actualResponse.PositionID, ShouldEqual, expectedResponse.PositionID)
				})
			})

			Convey("When creating an employee with invalid data", func() {
				// Create request payload with missing required fields
				reqInvalid := CreateRequest{
					Name:    "",
					Age:     0,
					Address: "",
					Phone:   "",
					Email:   "",
				}

				// Make the request and verify response
				var actualResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPost,
					"/employee",
					reqInvalid,
					&actualResponse,
					http.StatusBadRequest,
				)

				Convey("Then the response should indicate a validation error", func() {
					So(actualResponse["error"], ShouldNotBeEmpty)
				})
			})
		})
	})
}
