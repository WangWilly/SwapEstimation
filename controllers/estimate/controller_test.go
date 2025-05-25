package estimate

import (
	"testing"

	"github.com/WangWilly/swap-estimation/pkgs/testutils"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/mock/gomock"
)

////////////////////////////////////////////////////////////////////////////////

type testSuite struct {
	ethClient    *MockEthClient
	ethWssClient *MockEthWssClient

	controller *Controller
	testServer testutils.TestHttpServer
}

func testInit(t *testing.T, test func(*testSuite)) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ethClient := NewMockEthClient(ctrl)
	ethWssClient := NewMockEthWssClient(ctrl)
	cfg := Config{}
	if err := envconfig.Process(t.Context(), &cfg); err != nil {
		t.Fatal(err)
	}

	controller := NewController(cfg, ethClient, ethWssClient)
	testServer := testutils.NewTestHttpServer(controller)
	suite := &testSuite{
		ethClient:    ethClient,
		ethWssClient: ethWssClient,
		controller:   controller,
		testServer:   testServer,
	}

	test(suite)
}
