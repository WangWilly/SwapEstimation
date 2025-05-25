package ethwss

import (
	"testing"
	"time"

	"github.com/sethvargo/go-envconfig"
	"go.uber.org/mock/gomock"
)

////////////////////////////////////////////////////////////////////////////////

type testSuite struct {
	gethWssClient *MockGethWssClient

	client *client
}

func testInit(t *testing.T, test func(*testSuite)) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gethWssClient := NewMockGethWssClient(ctrl)

	cfg := Config{
		ListenPairPeriod: 2 * time.Minute,
	}
	if err := envconfig.Process(t.Context(), &cfg); err != nil {
		t.Fatal(err)
	}
	client := New(cfg, gethWssClient)

	ts := &testSuite{
		gethWssClient: gethWssClient,
		client:        client,
	}
	test(ts)
}
