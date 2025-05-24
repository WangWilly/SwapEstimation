package employee

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/models"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestPromote(t *testing.T) {
	testInit(t, func(s *testSuite) {
		Convey("Given an employee exists in the system", t, func() {
			// Setup test data
			employeeID := int64(123)
			nowTime := time.Date(2023, 6, 15, 12, 0, 0, 0, time.UTC)
			startDate := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC).Unix()

			// Position data for promotion
			newPosition := &models.EmployeePosition{
				ID:         456,
				EmployeeID: employeeID,
				Position:   "Senior Manager",
				Department: "Operations",
				Salary:     120000.00,
				StartDate:  time.Unix(startDate, 0),
			}

			// Promotion request
			req := PromoteRequest{
				Position:   newPosition.Position,
				Department: newPosition.Department,
				Salary:     newPosition.Salary,
				StartDate:  startDate,
			}

			Convey("When promoting the employee to a new position", func(c C) {
				// Set up expectations
				s.timeModule.EXPECT().Now().Return(nowTime)

				s.employeePositionRepo.EXPECT().
					Create(gomock.Any(), s.db, gomock.Any(), nowTime).
					DoAndReturn(func(_ interface{}, _ interface{}, position *models.EmployeePosition, _ time.Time) error {
						// Verify the position data
						c.So(position.EmployeeID, ShouldEqual, employeeID)
						c.So(position.Position, ShouldEqual, req.Position)
						c.So(position.Department, ShouldEqual, req.Department)
						c.So(position.Salary, ShouldEqual, req.Salary)
						c.So(position.StartDate, ShouldEqual, time.Unix(req.StartDate, 0))

						// Set the ID for the response
						position.ID = newPosition.ID
						return nil
					})

					// Expect cache to be deleted after successful promotion
				s.cacheManager.EXPECT().
					DeleteEmployeeDetailV1(gomock.Any(), employeeID).
					Return(nil)

				// Expected response
				expectedResponse := PromoteResponse{
					PositionID: newPosition.ID,
				}

				// Make the request and verify response
				var actualResponse PromoteResponse
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPost,
					"/promote/123",
					req,
					&actualResponse,
					http.StatusOK,
				)

				Convey("Then the response should contain the correct position ID", func() {
					So(actualResponse.PositionID, ShouldEqual, expectedResponse.PositionID)
				})
			})

			Convey("When providing an invalid employee ID", func() {
				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPost,
					"/promote/invalid-id",
					req,
					&errorResponse,
					http.StatusBadRequest,
				)

				Convey("Then the response should indicate invalid ID", func() {
					So(errorResponse["error"], ShouldEqual, "invalid id")
				})
			})

			Convey("When missing required fields in the request", func() {
				// Invalid request missing required fields
				invalidReq := map[string]interface{}{
					"position": "Senior Manager",
					// Missing department, salary, and start_date
				}

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPost,
					"/promote/123",
					invalidReq,
					&errorResponse,
					http.StatusBadRequest,
				)

				Convey("Then the response should indicate a validation error", func() {
					So(errorResponse["error"], ShouldNotBeEmpty)
				})
			})

			Convey("When creating the position record fails", func() {
				// Set up expectations for failure
				s.timeModule.EXPECT().Now().Return(nowTime)

				s.employeePositionRepo.EXPECT().
					Create(gomock.Any(), s.db, gomock.Any(), nowTime).
					Return(errors.New("database error"))

					// Cache deletion should not be called when database operation fails

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodPost,
					"/promote/123",
					req,
					&errorResponse,
					http.StatusInternalServerError,
				)

				Convey("Then the response should indicate a server error", func() {
					So(errorResponse["error"], ShouldEqual, "failed to create employee position")
				})
			})
		})
	})
}
