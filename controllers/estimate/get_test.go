package estimate

import (
	"context"
	"math/big"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/WangWilly/swap-estimation/pkgs/clients/eth"
	"github.com/WangWilly/swap-estimation/pkgs/clients/ethwss"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestGet(t *testing.T) {
	testInit(t, func(s *testSuite) {
		Convey("Given an estimate swap endpoint", t, func() {
			// Setup test data
			validPoolAddr := "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc"
			validSrcAddr := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2" // WETH
			validDstAddr := "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48" // USDC
			validAmount := "1000000000000000000"                         // 1 ETH

			mockReservePair := &eth.ReservePair{
				Reserve0: big.NewInt(200000000000), // 200,000 USDC (6 decimals)
				Reserve1: new(big.Int),             // 100 ETH (18 decimals)
			}
			mockReservePair.Reserve1.SetString("100000000000000000000", 10) // 100 ETH in wei

			mockEthWssReservePair := (*ethwss.ReservePair)(mockReservePair)

			expectedOutput := "1974316068"

			Convey("When making a valid estimation request", func() {
				// Set up expectations for the cache miss and eth client call
				s.ethWssClient.EXPECT().
					GetPair(gomock.Any(), validPoolAddr).
					Return(nil) // Cache miss

				s.ethClient.EXPECT().
					UniV2ReservePair(gomock.Any(), validPoolAddr).
					Return(mockReservePair, nil)

				s.ethWssClient.EXPECT().
					RegPair(gomock.Any(), validPoolAddr, gomock.Any()).
					Return(nil)

				// Make the request and verify response
				var actualOutput string
				resCode := s.testServer.MustDo(
					t,
					http.MethodGet,
					"/estimate?pool="+validPoolAddr+
						"&src="+validSrcAddr+
						"&dst="+validDstAddr+
						"&src_amount="+validAmount,
					nil,
					&actualOutput,
				)

				Convey("Then the response should be successful with the estimated amount", func() {
					So(resCode, ShouldEqual, http.StatusOK)
					So(actualOutput, ShouldEqual, expectedOutput)
				})
			})

			Convey("When making a request with cached pool data", func() {
				// Set up expectation for cache hit
				s.ethWssClient.EXPECT().
					GetPair(gomock.Any(), validPoolAddr).
					Return(mockEthWssReservePair) // Cache hit

				// No call to ethClient.UniV2ReservePair expected
				// No call to ethWssClient.RegPair expected

				// Make the request and verify response
				var actualOutput string
				resCode := s.testServer.MustDo(
					t,
					http.MethodGet,
					"/estimate?pool="+validPoolAddr+
						"&src="+validSrcAddr+
						"&dst="+validDstAddr+
						"&src_amount="+validAmount,
					nil,
					&actualOutput,
				)

				Convey("Then the response should be successful with the estimated amount from cache", func() {
					So(resCode, ShouldEqual, http.StatusOK)
					So(actualOutput, ShouldEqual, expectedOutput)
				})
			})

			Convey("When making a request with missing pool address", func() {
				// Make the request without pool parameter
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/estimate?src="+validSrcAddr+
						"&dst="+validDstAddr+
						"&src_amount="+validAmount,
					nil,
					&errorResponse,
					http.StatusBadRequest,
				)

				Convey("Then the response should indicate pool address is required", func() {
					So(errorResponse["error"], ShouldEqual, "invalid query parameters")
				})
			})

			Convey("When making a request with invalid pool address format", func() {
				invalidPoolAddr := "0xinvalid"

				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/estimate?pool="+invalidPoolAddr+
						"&src="+validSrcAddr+
						"&dst="+validDstAddr+
						"&src_amount="+validAmount,
					nil,
					&errorResponse,
					http.StatusBadRequest,
				)

				Convey("Then the response should indicate invalid pool address format", func() {
					So(errorResponse["error"], ShouldEqual, "invalid pool address format")
				})
			})

			Convey("When making a request with missing source token address", func() {
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/estimate?pool="+validPoolAddr+
						"&dst="+validDstAddr+
						"&src_amount="+validAmount,
					nil,
					&errorResponse,
					http.StatusBadRequest,
				)

				Convey("Then the response should indicate source token address is required", func() {
					So(errorResponse["error"], ShouldEqual, "invalid query parameters")
				})
			})

			Convey("When making a request with invalid source token address", func() {
				invalidSrcAddr := "0xinvalid"

				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/estimate?pool="+validPoolAddr+
						"&src="+invalidSrcAddr+
						"&dst="+validDstAddr+
						"&src_amount="+validAmount,
					nil,
					&errorResponse,
					http.StatusBadRequest,
				)

				Convey("Then the response should indicate invalid source token address format", func() {
					So(errorResponse["error"], ShouldEqual, "invalid source token address format")
				})
			})

			Convey("When making a request with missing destination token address", func() {
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/estimate?pool="+validPoolAddr+
						"&src="+validSrcAddr+
						"&src_amount="+validAmount,
					nil,
					&errorResponse,
					http.StatusBadRequest,
				)

				Convey("Then the response should indicate destination token address is required", func() {
					So(errorResponse["error"], ShouldEqual, "invalid query parameters")
				})
			})

			Convey("When making a request with invalid destination token address", func() {
				invalidDstAddr := "0xinvalid"

				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/estimate?pool="+validPoolAddr+
						"&src="+validSrcAddr+
						"&dst="+invalidDstAddr+
						"&src_amount="+validAmount,
					nil,
					&errorResponse,
					http.StatusBadRequest,
				)

				Convey("Then the response should indicate invalid destination token address format", func() {
					So(errorResponse["error"], ShouldEqual, "invalid destination token address format")
				})
			})

			Convey("When making a request with invalid pair address", func() {
				// Using valid addresses but they don't form a valid pair
				invalidPairPoolAddr := "0x1000000000000000000000000000000000000000"

				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/estimate?pool="+invalidPairPoolAddr+
						"&src="+validSrcAddr+
						"&dst="+validDstAddr+
						"&src_amount="+validAmount,
					nil,
					&errorResponse,
					http.StatusBadRequest,
				)

				Convey("Then the response should indicate invalid Uniswap V2 pair address", func() {
					So(errorResponse["error"], ShouldEqual, "invalid Uniswap V2 pair address")
				})
			})

			Convey("When making a request with missing source amount", func() {
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/estimate?pool="+validPoolAddr+
						"&src="+validSrcAddr+
						"&dst="+validDstAddr,
					nil,
					&errorResponse,
					http.StatusBadRequest,
				)

				Convey("Then the response should indicate source amount is required", func() {
					So(errorResponse["error"], ShouldEqual, "invalid query parameters")
				})
			})

			Convey("When making a request with invalid source amount", func() {
				invalidAmount := "not-a-number"

				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/estimate?pool="+validPoolAddr+
						"&src="+validSrcAddr+
						"&dst="+validDstAddr+
						"&src_amount="+invalidAmount,
					nil,
					&errorResponse,
					http.StatusBadRequest,
				)

				Convey("Then the response should indicate invalid source amount", func() {
					So(errorResponse["error"], ShouldEqual, "invalid source amount")
				})
			})

			Convey("When the eth client fails to retrieve reserves", func() {
				// Set up expectations for the cache miss and eth client failure
				s.ethWssClient.EXPECT().
					GetPair(gomock.Any(), validPoolAddr).
					Return(nil) // Cache miss

				// Set up expectation for eth client failure
				s.ethClient.EXPECT().
					UniV2ReservePair(gomock.Any(), validPoolAddr).
					Return(nil, context.DeadlineExceeded)

				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/estimate?pool="+validPoolAddr+
						"&src="+validSrcAddr+
						"&dst="+validDstAddr+
						"&src_amount="+validAmount,
					nil,
					&errorResponse,
					http.StatusInternalServerError,
				)

				Convey("Then the response should indicate failure to estimate output amount", func() {
					So(errorResponse["error"], ShouldEqual, "failed to get reserve pair")
				})
			})

			Convey("When multiple concurrent requests are made for the same pool", func() {
				// First request will be a cache miss
				s.ethWssClient.EXPECT().
					GetPair(gomock.Any(), validPoolAddr).
					Return(nil).
					Times(5) // First request is a cache miss

				// Set up expectation for the eth client call - should only be called ONCE
				s.ethClient.EXPECT().
					UniV2ReservePair(gomock.Any(), validPoolAddr).
					DoAndReturn(func(ctx context.Context, poolAddr string) (*eth.ReservePair, error) {
						// Simulate some processing time
						time.Sleep(100 * time.Millisecond)
						return mockReservePair, nil
					}).Times(1) // This is key - we expect only one call

				// Register the pair in cache
				s.ethWssClient.EXPECT().
					RegPair(gomock.Any(), validPoolAddr, gomock.Any()).
					Return(nil).
					Times(5) // Should only be called once after the first request

				// Number of concurrent requests to simulate
				concurrentRequests := 5
				var wg sync.WaitGroup
				wg.Add(concurrentRequests)

				// Channel to collect results
				resultChan := make(chan string, concurrentRequests)
				statusChan := make(chan int, concurrentRequests)

				// Make multiple "concurrent" requests
				for i := 0; i < concurrentRequests; i++ {
					go func() {
						defer wg.Done()
						var output string
						statusCode := s.testServer.MustDo(
							t,
							http.MethodGet,
							"/estimate?pool="+validPoolAddr+
								"&src="+validSrcAddr+
								"&dst="+validDstAddr+
								"&src_amount="+validAmount,
							nil,
							&output,
						)
						resultChan <- output
						statusChan <- statusCode
					}()
				}

				// Wait for all goroutines to complete
				wg.Wait()
				close(resultChan)
				close(statusChan)

				// Collect results
				results := make([]string, 0, concurrentRequests)
				statuses := make([]int, 0, concurrentRequests)

				for result := range resultChan {
					results = append(results, result)
				}

				for status := range statusChan {
					statuses = append(statuses, status)
				}

				Convey("Then all requests should receive successful responses", func() {
					for _, status := range statuses {
						So(status, ShouldEqual, http.StatusOK)
					}
				})

				Convey("Then all requests should receive the same output", func() {
					for _, result := range results {
						So(result, ShouldEqual, expectedOutput)
					}
				})

				// The key assertion is implicit in the Times(1) expectation above:
				// If the test passes, it means the ethClient.UniV2ReservePair was called exactly once
				// despite multiple concurrent requests
			})
		})
	})
}
