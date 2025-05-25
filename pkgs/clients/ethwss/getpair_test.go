package ethwss

import (
	"context"
	"math/big"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetPair(t *testing.T) {
	testInit(t, func(s *testSuite) {
		Convey("Given the GetPair function", t, func() {
			ctx := context.Background()
			pairAddr := "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc" // WETH-USDC pair

			Convey("When getting a pair that exists in cache", func() {
				// Create a reserve pair for testing
				reservePair := &ReservePair{
					Reserve0: big.NewInt(5000),
					Reserve1: big.NewInt(10000),
				}

				// Pre-populate the cache map
				s.client.reservePairCacheMap[pairAddr] = reservePair

				// Create a mock timer
				mockTimer := time.NewTimer(s.client.cfg.ListenPairPeriod)
				s.client.pairTimers[pairAddr] = mockTimer

				// Call the function
				result := s.client.GetPair(ctx, pairAddr)

				Convey("Then it should return the pair from cache", func() {
					So(result, ShouldNotBeNil)
					So(result.Reserve0.Cmp(reservePair.Reserve0), ShouldEqual, 0)
					So(result.Reserve1.Cmp(reservePair.Reserve1), ShouldEqual, 0)

					// Verify the timer was reset
					// Note: We can't easily test if timer was reset, but we can verify it exists
					timer, exists := s.client.pairTimers[pairAddr]
					So(exists, ShouldBeTrue)
					So(timer, ShouldEqual, mockTimer)
				})
			})

			Convey("When getting a pair that doesn't exist in cache", func() {
				// Call the function with a non-existent address
				delete(s.client.reservePairCacheMap, pairAddr) // Ensure it's not in cache
				result := s.client.GetPair(ctx, pairAddr)

				Convey("Then it should return nil", func() {
					So(result, ShouldBeNil)
				})
			})

			Convey("When getting a pair with existing cache but no timer", func() {
				// Create a reserve pair for testing
				reservePair := &ReservePair{
					Reserve0: big.NewInt(5000),
					Reserve1: big.NewInt(10000),
				}

				// Pre-populate the cache map but not the timer
				s.client.reservePairCacheMap[pairAddr] = reservePair

				// Call the function
				result := s.client.GetPair(ctx, pairAddr)

				Convey("Then it should return the pair from cache without error", func() {
					So(result, ShouldNotBeNil)
					So(result.Reserve0.Cmp(reservePair.Reserve0), ShouldEqual, 0)
					So(result.Reserve1.Cmp(reservePair.Reserve1), ShouldEqual, 0)
				})
			})
		})
	})
}
