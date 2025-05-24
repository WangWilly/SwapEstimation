package employeeinforepo

import (
	"testing"

	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"github.com/WangWilly/labs-hr-go/pkgs/testutils"
	"github.com/brianvoe/gofakeit/v6"
	"gorm.io/gorm"

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
		employeeInfo := models.DummyEmployeeInfo(faker)

		testutils.MustClearTable(t, db, models.EmployeeInfo{})

		// Create
		{
			Print("Create")

			err := repo.Create(ctx, db, employeeInfo)
			So(err, ShouldBeNil)
		}

		// Get
		{
			Print("Get")

			employeeInfoRes, err := repo.Get(ctx, db, employeeInfo.ID)
			So(err, ShouldBeNil)
			So(employeeInfoRes, ShouldNotBeNil)
			So(employeeInfoRes.ID, ShouldEqual, employeeInfo.ID)
			So(employeeInfoRes.Name, ShouldEqual, employeeInfo.Name)
			So(employeeInfoRes.Age, ShouldEqual, employeeInfo.Age)
			So(employeeInfoRes.Address, ShouldEqual, employeeInfo.Address)
			So(employeeInfoRes.Phone, ShouldEqual, employeeInfo.Phone)
			So(employeeInfoRes.Email, ShouldEqual, employeeInfo.Email)
			So(employeeInfoRes.CreatedAt, ShouldHappenOnOrAfter, employeeInfo.CreatedAt)
			So(employeeInfoRes.UpdatedAt, ShouldHappenOnOrAfter, employeeInfo.UpdatedAt)
			So(employeeInfoRes.DeleteAt, ShouldResemble, gorm.DeletedAt{})
		}

		// Save
		{
			Print("Save")

			employeeInfo.Name = faker.Name()
			employeeInfo.Age = faker.Number(20, 50)
			employeeInfo.Address = faker.StreetName()
			employeeInfo.Phone = faker.Phone()
			employeeInfo.Email = faker.Email()

			err := repo.Save(ctx, db, employeeInfo)
			So(err, ShouldBeNil)

			employeeInfoRes, err := repo.Get(ctx, db, employeeInfo.ID)
			So(err, ShouldBeNil)
			So(employeeInfoRes, ShouldNotBeNil)
			So(employeeInfoRes.ID, ShouldEqual, employeeInfo.ID)
			So(employeeInfoRes.Name, ShouldEqual, employeeInfo.Name)
			So(employeeInfoRes.Age, ShouldEqual, employeeInfo.Age)
			So(employeeInfoRes.Address, ShouldEqual, employeeInfo.Address)
			So(employeeInfoRes.Phone, ShouldEqual, employeeInfo.Phone)
			So(employeeInfoRes.Email, ShouldEqual, employeeInfo.Email)
			So(employeeInfoRes.CreatedAt, ShouldHappenOnOrAfter, employeeInfo.CreatedAt)
			So(employeeInfoRes.UpdatedAt, ShouldHappenOnOrAfter, employeeInfo.UpdatedAt)
			So(employeeInfoRes.DeleteAt, ShouldResemble, gorm.DeletedAt{})
		}
	})
}
