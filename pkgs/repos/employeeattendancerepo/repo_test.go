package employeeattendancerepo

import (
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
		employeeAttendance := models.DummyEmployeeAttendance(faker)

		testutils.MustClearTable(t, db, models.EmployeeAttendance{})

		// Create
		{
			Print("Create")

			employeeAttendanceRes, err := repo.CreateForClockIn(ctx, db, employeeAttendance.EmployeeID, employeeAttendance.PositionID, employeeAttendance.ClockIn)
			So(err, ShouldBeNil)
			So(employeeAttendanceRes, ShouldNotBeNil)
			So(employeeAttendanceRes.EmployeeID, ShouldEqual, employeeAttendance.EmployeeID)
			So(employeeAttendanceRes.PositionID, ShouldEqual, employeeAttendance.PositionID)
			So(employeeAttendanceRes.ClockIn, ShouldHappenWithin, time.Second, employeeAttendance.ClockIn)
			So(employeeAttendanceRes.ClockOut, ShouldHappenWithin, time.Second, employeeAttendance.ClockIn)
			So(employeeAttendanceRes.CreatedAt.IsZero(), ShouldBeFalse)
			So(employeeAttendanceRes.UpdatedAt.IsZero(), ShouldBeFalse)
			So(employeeAttendanceRes.ID, ShouldNotEqual, 0)

			employeeAttendance.ID = employeeAttendanceRes.ID
		}

		// Get
		{
			Print("Get")

			employeeAttendanceRes, err := repo.Last(ctx, db, employeeAttendance.EmployeeID)
			So(err, ShouldBeNil)
			So(employeeAttendanceRes, ShouldNotBeNil)
			So(employeeAttendanceRes.ID, ShouldEqual, employeeAttendance.ID)
			So(employeeAttendanceRes.EmployeeID, ShouldEqual, employeeAttendance.EmployeeID)
			So(employeeAttendanceRes.PositionID, ShouldEqual, employeeAttendance.PositionID)
			So(employeeAttendanceRes.ClockIn, ShouldHappenWithin, time.Second, employeeAttendance.ClockIn)
			So(employeeAttendanceRes.ClockOut, ShouldHappenWithin, time.Second, employeeAttendance.ClockIn)
			So(employeeAttendanceRes.CreatedAt.IsZero(), ShouldBeFalse)
			So(employeeAttendanceRes.UpdatedAt.IsZero(), ShouldBeFalse)
		}

		// Update
		{
			Print("Update")

			employeeAttendanceRes, err := repo.UpdateForClockOut(ctx, db, employeeAttendance.ID, employeeAttendance.ClockOut)
			So(err, ShouldBeNil)
			So(employeeAttendanceRes, ShouldNotBeNil)
			So(employeeAttendanceRes.ID, ShouldEqual, employeeAttendance.ID)
			So(employeeAttendanceRes.EmployeeID, ShouldEqual, employeeAttendance.EmployeeID)
			So(employeeAttendanceRes.PositionID, ShouldEqual, employeeAttendance.PositionID)
			So(employeeAttendanceRes.ClockIn, ShouldHappenWithin, time.Second, employeeAttendance.ClockIn)
			So(employeeAttendanceRes.ClockOut, ShouldHappenWithin, time.Second, employeeAttendance.ClockOut)
			So(employeeAttendanceRes.CreatedAt.IsZero(), ShouldBeFalse)
			So(employeeAttendanceRes.UpdatedAt.IsZero(), ShouldBeFalse)
		}

		{
			Print("Update - Already clocked out")

			employeeAttendanceRes, err := repo.UpdateForClockOut(ctx, db, employeeAttendance.ID, employeeAttendance.ClockOut)
			So(err, ShouldNotBeNil)
			So(employeeAttendanceRes, ShouldBeNil)
			So(err.Error(), ShouldEqual, "attendance record already clocked out")
		}
	})
}
