package ethwss

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestRegPair(t *testing.T) {
	testInit(t, func(s *testSuite) {
		Convey("Given the RegPair function", t, func() {
			ctx := context.Background()
			pairAddr := "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc"                             // WETH-USDC pair
			syncEventSig := "0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1" // Sync event signature

			initPair := &ReservePair{
				Reserve0: big.NewInt(5000),
				Reserve1: big.NewInt(10000),
			}

			Convey("When registering a new pair", func(c C) {
				// Mock the subscription
				mockSub := new(mockSubscription)

				// Mock FilterLogs call
				s.gethWssClient.EXPECT().
					SubscribeFilterLogs(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
						c.So(query.Addresses, ShouldContain, common.HexToAddress(pairAddr))
						c.So(query.Topics, ShouldHaveLength, 1)
						c.So(query.Topics[0], ShouldHaveLength, 1)
						c.So(query.Topics[0][0].Hex(), ShouldEqual, syncEventSig)
						return mockSub, nil
					})

				// Call the function
				err := s.client.RegPair(ctx, pairAddr, initPair)

				Convey("Then it should register the pair without error", func() {
					So(err, ShouldBeNil)

					// Verify the pair is in the cache map
					pair, exists := s.client.reservePairCacheMap[pairAddr]
					So(exists, ShouldBeTrue)
					So(pair.Reserve0.Cmp(initPair.Reserve0), ShouldEqual, 0)
					So(pair.Reserve1.Cmp(initPair.Reserve1), ShouldEqual, 0)

					// Verify the subscription is saved
					sub, exists := s.client.pairSubscriptions[pairAddr]
					So(exists, ShouldBeTrue)
					So(sub, ShouldEqual, mockSub)

					// Verify the timer is created
					timer, exists := s.client.pairTimers[pairAddr]
					So(exists, ShouldBeTrue)
					So(timer, ShouldNotBeNil)
				})
			})

			Convey("When registering a pair that's already registered", func() {
				// Pre-populate the cache map
				s.client.reservePairCacheMap[pairAddr] = initPair

				// Call the function
				err := s.client.RegPair(ctx, pairAddr, initPair)

				Convey("Then it should skip registration without error", func() {
					So(err, ShouldBeNil)
				})
			})

			Convey("When registration for a pair is already in progress", func() {
				// Mark the pair as being registered
				s.client.addressLock.Lock()
				s.client.registeringPairs[pairAddr] = true
				delete(s.client.reservePairCacheMap, pairAddr) // Ensure it's not in cache
				s.client.addressLock.Unlock()

				// Call the function
				err := s.client.RegPair(ctx, pairAddr, initPair)

				Convey("Then it should skip without error", func() {
					So(err, ShouldBeNil)

					// Verify the cache map was not modified
					_, exists := s.client.reservePairCacheMap[pairAddr]
					So(exists, ShouldBeFalse)
				})
			})

			Convey("When subscription fails", func() {
				s.client.addressLock.Lock()
				s.client.registeringPairs[pairAddr] = false
				s.client.addressLock.Unlock()
				// Mock a failed subscription
				s.gethWssClient.EXPECT().
					SubscribeFilterLogs(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("subscription failed"))

				// Call the function
				err := s.client.RegPair(ctx, pairAddr, initPair)

				Convey("Then it should return an error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "subscription failed")

					// Verify the registration flag was cleared
					s.client.addressLock.Lock()
					registering := s.client.registeringPairs[pairAddr]
					s.client.addressLock.Unlock()
					So(registering, ShouldBeFalse)
				})
			})
		})
	})
}

// Mock implementation of ethereum.Subscription
type mockSubscription struct{}

func (m *mockSubscription) Unsubscribe()      {}
func (m *mockSubscription) Err() <-chan error { return make(<-chan error) }
