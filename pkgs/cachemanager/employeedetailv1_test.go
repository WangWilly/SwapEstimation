package cachemanager

import (
	"testing"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/dtos"
	. "github.com/smartystreets/goconvey/convey"
)

////////////////////////////////////////////////////////////////////////////////

func TestEmployeeDetailV1(t *testing.T) {
	testInit(t, func(s *testSuite) {
		Convey("Given an employee cache manager", t, func() {
			ctx := t.Context()
			employeeID := int64(123)

			// Create test data
			employeeData := dtos.EmployeeV1Response{
				EmployeeID: employeeID,
				Name:       "John Doe",
				Age:        30,
				Phone:      "555-1234",
				Email:      "john.doe@example.com",
				Address:    "123 Main St",
				CreatedAt:  "2023-05-15T10:00:00Z",
				UpdatedAt:  "2023-05-15T10:00:00Z",
				PositionID: 456,
				Position:   "Developer",
				Department: "Engineering",
				Salary:     75000.00,
				StartDate:  "2023-01-01T00:00:00Z",
			}

			// Clean up before testing to ensure consistent state
			_ = s.manager.DeleteEmployeeDetailV1(ctx, employeeID)

			Convey("When setting employee detail cache", func() {
				err := s.manager.SetEmployeeDetailV1(ctx, employeeID, employeeData, time.Minute*15)

				Convey("Then no error should occur", func() {
					So(err, ShouldBeNil)
				})

				Convey("And when getting the cached employee detail", func() {
					cachedData, err := s.manager.GetEmployeeDetailV1(ctx, employeeID)

					Convey("Then no error should occur", func() {
						So(err, ShouldBeNil)
						So(cachedData, ShouldNotBeNil)
					})

					Convey("And the cached data should match the original data", func() {
						So(cachedData.EmployeeID, ShouldEqual, employeeData.EmployeeID)
						So(cachedData.Name, ShouldEqual, employeeData.Name)
						So(cachedData.Age, ShouldEqual, employeeData.Age)
						So(cachedData.Phone, ShouldEqual, employeeData.Phone)
						So(cachedData.Email, ShouldEqual, employeeData.Email)
						So(cachedData.Address, ShouldEqual, employeeData.Address)
						So(cachedData.PositionID, ShouldEqual, employeeData.PositionID)
						So(cachedData.Position, ShouldEqual, employeeData.Position)
						So(cachedData.Department, ShouldEqual, employeeData.Department)
						So(cachedData.Salary, ShouldEqual, employeeData.Salary)
						So(cachedData.StartDate, ShouldEqual, employeeData.StartDate)
					})
				})
			})

			Convey("When updating existing employee detail cache", func() {
				// First set the initial data
				err := s.manager.SetEmployeeDetailV1(ctx, employeeID, employeeData, time.Minute*15)
				So(err, ShouldBeNil)

				// Update the data
				updatedData := employeeData
				updatedData.Name = "Jane Smith"
				updatedData.Age = 32
				updatedData.Position = "Senior Developer"
				updatedData.Salary = 85000.00

				// Set the updated data
				err = s.manager.SetEmployeeDetailV1(ctx, employeeID, updatedData, time.Minute*15)

				Convey("Then no error should occur", func() {
					So(err, ShouldBeNil)
				})

				Convey("And when getting the updated cached employee detail", func() {
					cachedData, err := s.manager.GetEmployeeDetailV1(ctx, employeeID)

					Convey("Then the cached data should contain the updates", func() {
						So(err, ShouldBeNil)
						So(cachedData, ShouldNotBeNil)
						So(cachedData.Name, ShouldEqual, updatedData.Name)
						So(cachedData.Age, ShouldEqual, updatedData.Age)
						So(cachedData.Position, ShouldEqual, updatedData.Position)
						So(cachedData.Salary, ShouldEqual, updatedData.Salary)
					})
				})
			})

			Convey("When getting a non-existent employee detail", func() {
				nonExistentID := int64(999)
				// Make sure it doesn't exist
				_ = s.manager.DeleteEmployeeDetailV1(ctx, nonExistentID)

				cachedData, err := s.manager.GetEmployeeDetailV1(ctx, nonExistentID)

				Convey("Then an error should not occur", func() {
					So(err, ShouldBeNil)
				})

				Convey("And the cached data should be nil", func() {
					So(cachedData, ShouldBeNil)
				})
			})

			Convey("When deleting an employee detail cache", func() {
				// First set the data
				err := s.manager.SetEmployeeDetailV1(ctx, employeeID, employeeData, time.Minute*15)
				So(err, ShouldBeNil)

				// Then delete it
				err = s.manager.DeleteEmployeeDetailV1(ctx, employeeID)

				Convey("Then an error should not occur", func() {
					So(err, ShouldBeNil)
				})

				Convey("And when trying to get the deleted cache", func() {
					cachedData, err := s.manager.GetEmployeeDetailV1(ctx, employeeID)

					Convey("Then an error should not occur", func() {
						So(err, ShouldBeNil)
					})

					Convey("And the cached data should be nil", func() {
						So(cachedData, ShouldBeNil)
					})
				})
			})

			Convey("When deleting a non-existent employee detail cache", func() {
				nonExistentID := int64(999)
				// Make sure it doesn't exist
				_ = s.manager.DeleteEmployeeDetailV1(ctx, nonExistentID)

				// Try to delete it again
				err := s.manager.DeleteEmployeeDetailV1(ctx, nonExistentID)

				Convey("Then no error should occur", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	})
}
