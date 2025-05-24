package eth

import (
	"testing"

	"github.com/sethvargo/go-envconfig"
	"go.uber.org/mock/gomock"
)

////////////////////////////////////////////////////////////////////////////////

type testSuite struct {
	gethClient *MockGethClient

	client *client
}

func testInit(t *testing.T, test func(*testSuite)) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gethClient := NewMockGethClient(ctrl)

	cfg := Config{}
	if err := envconfig.Process(t.Context(), &cfg); err != nil {
		t.Fatal(err)
	}
	client := New(cfg, gethClient)

	ts := &testSuite{
		gethClient: gethClient,
		client:     client,
	}
	test(ts)
}
