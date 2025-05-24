package employeepositionrepo

import (
	"errors"
	"testing"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"github.com/WangWilly/labs-hr-go/pkgs/testutils"
	"github.com/brianvoe/gofakeit/v6"

	. "github.com/smartystreets/goconvey/convey"
)

////////////////////////////////////////////////////////////////////////////////

func TestMain(m *testing.M) {
	testutils.BeforeTestDb(m)
}

////////////////////////////////////////////////////////////////////////////////

func TestRepo_CRUD(t *testing.T) {
	Convey("TestRepo_CRUD", t, func() {
		// Setup
		ctx := t.Context()
		db := testutils.GetDB().DB
		repo := New()
		faker := gofakeit.New(0)

		// Prepare test data
		employeePosition := models.DummyEmployeePosition(faker)

		testutils.MustClearTable(t, db, models.EmployeePosition{})

		// Create
		{
			Print("Create")

			err := repo.Create(ctx, db, employeePosition, time.Now())
			So(err, ShouldBeNil)
		}

		// Get
		{
			Print("Get")

			employeePositionRes, err := repo.Get(ctx, db, employeePosition.ID)
			So(err, ShouldBeNil)
			So(employeePositionRes, ShouldNotBeNil)
			So(employeePositionRes.ID, ShouldEqual, employeePosition.ID)
			So(employeePositionRes.EmployeeID, ShouldEqual, employeePosition.EmployeeID)
			So(employeePositionRes.Position, ShouldEqual, employeePosition.Position)
			So(employeePositionRes.Department, ShouldEqual, employeePosition.Department)
			So(employeePositionRes.Salary, ShouldEqual, employeePosition.Salary)
			So(employeePositionRes.StartDate, ShouldEqual, employeePosition.StartDate)
			So(employeePositionRes.CreatedAt, ShouldHappenOnOrAfter, employeePosition.CreatedAt)
		}

		// GetCurrentByEmployeeID
		{
			Print("GetCurrentByEmployeeID")

			employeePositionRes, err := repo.GetCurrentByEmployeeID(ctx, db, employeePosition.EmployeeID, employeePosition.StartDate)
			So(err, ShouldBeNil)
			So(employeePositionRes, ShouldNotBeNil)
			So(employeePositionRes.ID, ShouldEqual, employeePosition.ID)
			So(employeePositionRes.EmployeeID, ShouldEqual, employeePosition.EmployeeID)
			So(employeePositionRes.Position, ShouldEqual, employeePosition.Position)
			So(employeePositionRes.Department, ShouldEqual, employeePosition.Department)
			So(employeePositionRes.Salary, ShouldEqual, employeePosition.Salary)
			So(employeePositionRes.StartDate, ShouldEqual, employeePosition.StartDate)
			So(employeePositionRes.CreatedAt, ShouldHappenOnOrAfter, employeePosition.CreatedAt)
		}

		// Create another employee position before the first one
		{
			Print("Create another employee position before the first one")

			employeePosition2 := models.DummyEmployeePosition(faker)
			employeePosition2.EmployeeID = employeePosition.EmployeeID
			employeePosition2.StartDate = employeePosition.StartDate.AddDate(0, 0, -1)
			err := repo.Create(ctx, db, employeePosition2, time.Now())
			So(err, ShouldResemble, errors.New("start_date must be greater than the current position's start_date"))
		}

		// Create another employee position after the first one
		// GetCurrentByEmployeeID
		{
			Print("Create another employee position after the first one")

			employeePosition2 := models.DummyEmployeePosition(faker)
			employeePosition2.EmployeeID = employeePosition.EmployeeID
			employeePosition2.StartDate = employeePosition.StartDate.AddDate(0, 0, 1)
			err := repo.Create(ctx, db, employeePosition2, time.Now())
			So(err, ShouldBeNil)

			Print("GetCurrentByEmployeeID")

			employeePositionRes, err := repo.GetCurrentByEmployeeID(ctx, db, employeePosition.EmployeeID, employeePosition.StartDate.AddDate(0, 0, 1))
			So(err, ShouldBeNil)
			So(employeePositionRes, ShouldNotBeNil)
			So(employeePositionRes.ID, ShouldEqual, employeePosition2.ID)
			So(employeePositionRes.EmployeeID, ShouldEqual, employeePosition2.EmployeeID)
			So(employeePositionRes.Position, ShouldEqual, employeePosition2.Position)
			So(employeePositionRes.Department, ShouldEqual, employeePosition2.Department)
			So(employeePositionRes.Salary, ShouldEqual, employeePosition2.Salary)
			So(employeePositionRes.StartDate, ShouldEqual, employeePosition2.StartDate)
			So(employeePositionRes.CreatedAt, ShouldHappenOnOrAfter, employeePosition2.CreatedAt)
		}
	})
}
