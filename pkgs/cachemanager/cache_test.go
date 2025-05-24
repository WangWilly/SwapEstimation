package cachemanager

import (
	"testing"

	"github.com/WangWilly/labs-hr-go/pkgs/testutils"
	"go.uber.org/mock/gomock"
)

////////////////////////////////////////////////////////////////////////////////

func TestMain(m *testing.M) {
	testutils.BeforeTestDbRedis(m)
}

////////////////////////////////////////////////////////////////////////////////

type testSuite struct {
	manager *manager
}

func testInit(t *testing.T, test func(*testSuite)) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	manager := New(testutils.GetRedis().RedisClient)

	test(&testSuite{
		manager: manager,
	})
}
