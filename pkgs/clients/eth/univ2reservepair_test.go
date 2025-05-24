package eth

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

func TestUniV2ReservePair(t *testing.T) {
	testInit(t, func(s *testSuite) {
		Convey("Given the UniV2ReservePair function", t, func() {
			ctx := context.Background()
			pairAddr := "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc"                             // WETH-USDC pair
			syncEventSig := "0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1" // Sync event signature

			// Set block range size for testing
			s.client.cfg.BlockRangeSize = 100

			Convey("When querying for pool reserves and logs are found in the first block range", func(c C) {
				// Mock responses
				latestBlock := uint64(15000000)

				// Setup the expected filter query
				expectedQuery := ethereum.FilterQuery{
					FromBlock: big.NewInt(14999900), // latestBlock - BlockRangeSize
					ToBlock:   big.NewInt(15000000), // latestBlock
					Addresses: []common.Address{common.HexToAddress(pairAddr)},
					Topics:    [][]common.Hash{{common.HexToHash(syncEventSig)}},
				}

				// Mock event logs
				mockLogs := []types.Log{
					{
						Address:     common.HexToAddress(pairAddr),
						Topics:      []common.Hash{common.HexToHash(syncEventSig)},
						Data:        []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x13, 0x88, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x27, 0x10},
						BlockNumber: 15000000 - 10,
						TxHash:      common.HexToHash("0x123"),
					},
					{
						Address:     common.HexToAddress(pairAddr),
						Topics:      []common.Hash{common.HexToHash(syncEventSig)},
						Data:        []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x13, 0x89, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x27, 0x11},
						BlockNumber: 15000000 - 5,
						TxHash:      common.HexToHash("0x456"),
					},
				}

				// Set up expectations
				s.gethClient.EXPECT().
					BlockNumber(gomock.Any()).
					Return(latestBlock, nil)

				s.gethClient.EXPECT().
					FilterLogs(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
						// Verify the query matches the expected one
						c.So(query.FromBlock.Int64(), ShouldEqual, expectedQuery.FromBlock.Int64())
						c.So(query.ToBlock.Int64(), ShouldEqual, expectedQuery.ToBlock.Int64())
						c.So(query.Addresses, ShouldContain, common.HexToAddress(pairAddr))
						c.So(query.Topics, ShouldHaveLength, 1)
						c.So(query.Topics[0], ShouldHaveLength, 1)
						c.So(query.Topics[0][0], ShouldEqual, common.HexToHash(syncEventSig))
						return mockLogs, nil
					})

				// Call the function
				result, err := s.client.UniV2ReservePair(ctx, pairAddr)

				Convey("Then it should return the latest reserve values without error", func() {
					So(err, ShouldBeNil)
					So(result, ShouldNotBeNil)

					// Expected values from the mock data (latest log)
					// Data contains reserve0 = 5001 (0x1389) and reserve1 = 10001 (0x2711)
					So(result.Reserve0.Cmp(big.NewInt(5001)), ShouldEqual, 0)
					So(result.Reserve1.Cmp(big.NewInt(10001)), ShouldEqual, 0)
				})
			})

			Convey("When logs are found after searching through multiple block ranges", func(c C) {
				// Mock responses
				latestBlock := uint64(15000000)

				// First query finds no logs
				firstQuery := ethereum.FilterQuery{
					FromBlock: big.NewInt(14999900), // latestBlock - BlockRangeSize
					ToBlock:   big.NewInt(15000000), // latestBlock
					Addresses: []common.Address{common.HexToAddress(pairAddr)},
					Topics:    [][]common.Hash{{common.HexToHash(syncEventSig)}},
				}

				// Second query finds logs
				secondQuery := ethereum.FilterQuery{
					FromBlock: big.NewInt(14999800), // latestBlock - 2*BlockRangeSize
					ToBlock:   big.NewInt(14999900), // latestBlock - BlockRangeSize
					Addresses: []common.Address{common.HexToAddress(pairAddr)},
					Topics:    [][]common.Hash{{common.HexToHash(syncEventSig)}},
				}

				// Mock event logs for second query
				mockLogs := []types.Log{
					{
						Address:     common.HexToAddress(pairAddr),
						Topics:      []common.Hash{common.HexToHash(syncEventSig)},
						Data:        []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xD, 0x48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0x86, 0xA0},
						BlockNumber: 14999850,
						TxHash:      common.HexToHash("0x789"),
					},
				}

				// Set up expectations
				s.gethClient.EXPECT().
					BlockNumber(gomock.Any()).
					Return(latestBlock, nil)

				s.gethClient.EXPECT().
					FilterLogs(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
						c.So(query.FromBlock.Int64(), ShouldEqual, firstQuery.FromBlock.Int64())
						c.So(query.ToBlock.Int64(), ShouldEqual, firstQuery.ToBlock.Int64())
						c.So(query.Addresses, ShouldContain, common.HexToAddress(pairAddr))
						c.So(query.Topics, ShouldHaveLength, 1)
						c.So(query.Topics[0], ShouldHaveLength, 1)
						c.So(query.Topics[0][0], ShouldEqual, common.HexToHash(syncEventSig))
						return []types.Log{}, nil // No logs found in first query
					})

				s.gethClient.EXPECT().
					FilterLogs(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
						c.So(query.FromBlock.Int64(), ShouldEqual, secondQuery.FromBlock.Int64())
						c.So(query.ToBlock.Int64(), ShouldEqual, secondQuery.ToBlock.Int64())
						c.So(query.Addresses, ShouldContain, common.HexToAddress(pairAddr))
						c.So(query.Topics, ShouldHaveLength, 1)
						c.So(query.Topics[0], ShouldHaveLength, 1)
						c.So(query.Topics[0][0], ShouldEqual, common.HexToHash(syncEventSig))
						return mockLogs, nil // Logs found in second query
					})

				// Call the function
				result, err := s.client.UniV2ReservePair(ctx, pairAddr)

				Convey("Then it should return the reserve values from the found logs", func() {
					So(err, ShouldBeNil)
					So(result, ShouldNotBeNil)

					// Expected values from the mock data
					// Data contains reserve0 = 3400 (0xD48) and reserve1 = 100000 (0x186A0)
					So(result.Reserve0.Cmp(big.NewInt(3400)), ShouldEqual, 0)
					So(result.Reserve1.Cmp(big.NewInt(100000)), ShouldEqual, 0)
				})
			})

			Convey("When no logs are found in any block range", func(c C) {
				// Mock responses
				latestBlock := s.client.cfg.BlockRangeSize * 3 // Ensure we have enough blocks to search

				// Set up expectations for many queries, all returning no logs
				s.gethClient.EXPECT().
					BlockNumber(gomock.Any()).
					Return(latestBlock, nil)

				// Need to set up multiple FilterLogs calls that will all return empty logs
				// This is a simplification - the real implementation would make many calls
				// until it reaches block 0 or finds logs
				blockRangeSize := s.client.cfg.BlockRangeSize
				maxQueries := 3 // Limit the number of queries for the test

				for i := 0; i < maxQueries; i++ {
					fromBlock := int64(latestBlock) - int64(uint64(i+1)*blockRangeSize)
					if fromBlock < 0 {
						fromBlock = 0
					}
					toBlock := int64(latestBlock) - int64(uint64(i)*blockRangeSize)

					s.gethClient.EXPECT().
						FilterLogs(gomock.Any(), gomock.Any()).
						DoAndReturn(func(_ context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
							c.So(query.FromBlock.Int64(), ShouldEqual, fromBlock)
							c.So(query.ToBlock.Int64(), ShouldEqual, toBlock)
							return []types.Log{}, nil
						})
				}

				// Call the function
				result, err := s.client.UniV2ReservePair(ctx, pairAddr)

				Convey("Then it should return an error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, "no Sync events found")
					So(result, ShouldBeNil)
				})
			})

			Convey("When there is an error getting the latest block number", func() {
				// Set up expectations
				s.gethClient.EXPECT().
					BlockNumber(gomock.Any()).
					Return(uint64(0), errors.New("blockchain connection error"))

				// Call the function
				result, err := s.client.UniV2ReservePair(ctx, pairAddr)

				Convey("Then it should return the error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, "failed to get latest block")
					So(result, ShouldBeNil)
				})
			})

			Convey("When there is an error filtering logs", func() {
				// Mock responses
				latestBlock := s.client.cfg.BlockRangeSize

				// Set up expectations
				s.gethClient.EXPECT().
					BlockNumber(gomock.Any()).
					Return(latestBlock, nil)

				s.gethClient.EXPECT().
					FilterLogs(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("filter logs error"))

				// Call the function
				result, err := s.client.UniV2ReservePair(ctx, pairAddr)

				Convey("Then it should continue searching in earlier blocks", func() {
					// In the actual implementation, it would continue searching
					// For the test, we're just verifying it doesn't immediately return an error
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, "no Sync events found")
					So(result, ShouldBeNil)
				})
			})

			Convey("When there is an error unpacking log data", func() {
				// Mock responses
				latestBlock := uint64(15000000)

				// Mock event logs with invalid data (too short)
				mockLogs := []types.Log{
					{
						Address:     common.HexToAddress(pairAddr),
						Topics:      []common.Hash{common.HexToHash(syncEventSig)},
						Data:        []byte{0x01, 0x02}, // Too short to be valid Sync event data
						BlockNumber: 15000000 - 10,
						TxHash:      common.HexToHash("0x123"),
					},
				}

				// Set up expectations
				s.gethClient.EXPECT().
					BlockNumber(gomock.Any()).
					Return(latestBlock, nil)

				s.gethClient.EXPECT().
					FilterLogs(gomock.Any(), gomock.Any()).
					Return(mockLogs, nil)

				// Call the function
				result, err := s.client.UniV2ReservePair(ctx, pairAddr)

				Convey("Then it should return an error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, "failed to unpack log data")
					So(result, ShouldBeNil)
				})
			})
		})
	})
}
