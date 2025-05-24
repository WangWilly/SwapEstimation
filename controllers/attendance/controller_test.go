package attendance

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/WangWilly/labs-hr-go/pkgs/testutils"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/sethvargo/go-envconfig"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

////////////////////////////////////////////////////////////////////////////////

type testSuite struct {
	db     *gorm.DB
	mockDB sqlmock.Sqlmock

	timeModule             *MockTimeModule
	employeePositionRepo   *MockEmployeePositionRepo
	employeeAttendanceRepo *MockEmployeeAttendanceRepo
	cacheManager           *MockCacheManager

	controller *Controller
	testServer testutils.TestHttpServer
	faker      *gofakeit.Faker
}

func testInit(t *testing.T, test func(*testSuite)) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gormDB, mockDB := testutils.GetMockDB(t)

	timeModule := NewMockTimeModule(ctrl)
	employeePositionRepo := NewMockEmployeePositionRepo(ctrl)
	employeeAttendanceRepo := NewMockEmployeeAttendanceRepo(ctrl)
	cacheManager := NewMockCacheManager(ctrl)

	cfg := Config{}
	if err := envconfig.Process(t.Context(), &cfg); err != nil {
		t.Fatal(err)
	}
	controller := NewController(
		cfg,
		gormDB,
		timeModule,
		employeePositionRepo,
		employeeAttendanceRepo,
		cacheManager,
	)
	testServer := testutils.NewTestHttpServer(controller)
	faker := gofakeit.New(0)
	suite := &testSuite{
		db:                     gormDB,
		mockDB:                 mockDB,
		timeModule:             timeModule,
		employeePositionRepo:   employeePositionRepo,
		employeeAttendanceRepo: employeeAttendanceRepo,
		cacheManager:           cacheManager,
		controller:             controller,
		testServer:             testServer,
		faker:                  faker,
	}

	test(suite)
}
