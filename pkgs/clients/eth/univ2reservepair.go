package eth

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

////////////////////////////////////////////////////////////////////////////////

type ReservePair struct {
	Reserve0 *big.Int
	Reserve1 *big.Int
}

////////////////////////////////////////////////////////////////////////////////

const uniswapV2PairABI = `[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint112","name":"reserve0","type":"uint112"},{"indexed":false,"internalType":"uint112","name":"reserve1","type":"uint112"}],"name":"Sync","type":"event"}]`

////////////////////////////////////////////////////////////////////////////////

func (c *client) UniV2ReservePair(
	ctx context.Context,
	pairAddrStr string,
) (*ReservePair, error) {
	logger := log.Ctx(ctx)
	logger.Debug().
		Str("pair_address", pairAddrStr).
		Msg("Estimating Uniswap V2 output amount")

	pairAddress := common.HexToAddress(pairAddrStr)
	parsedABI, err := abi.JSON(strings.NewReader(uniswapV2PairABI))
	if err != nil {
		logger.Error().Err(err).Msg("Failed to parse Uniswap V2 Pair ABI")
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	// Get latest block number
	latestBlock, err := c.gethClient.BlockNumber(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get latest block number")
		return nil, fmt.Errorf("failed to get latest block: %v", err)
	}

	var logs []types.Log
	var foundLogs bool

	// Start from the latest block and search backward in chunks
	for toBlock := latestBlock; toBlock > 0; toBlock -= c.cfg.BlockRangeSize {
		fromBlock := toBlock - c.cfg.BlockRangeSize
		if fromBlock > toBlock {
			fromBlock = 0
		}

		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(fromBlock)),
			ToBlock:   big.NewInt(int64(toBlock)),
			Addresses: []common.Address{pairAddress},
			Topics:    [][]common.Hash{{parsedABI.Events["Sync"].ID}},
		}

		chunkLogs, err := c.gethClient.FilterLogs(ctx, query)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to filter logs from block %d to %d", fromBlock, toBlock)
			continue
		}

		if len(chunkLogs) > 0 {
			logs = chunkLogs
			foundLogs = true
			break
		}
	}

	if !foundLogs {
		logger.Warn().Msg("No Sync events found in any block range")
		return nil, fmt.Errorf("no Sync events found in any block range")
	}

	latestLog := logs[len(logs)-1]
	var reserve0, reserve1 *big.Int
	if err := parsedABI.UnpackIntoInterface(&[]any{&reserve0, &reserve1}, "Sync", latestLog.Data); err != nil {
		logger.Error().Err(err).Msg("Failed to unpack log data")
		return nil, fmt.Errorf("failed to unpack log data: %v", err)
	}

	// Optional debug info
	logger.Debug().
		Str("pair_address", pairAddress.Hex()).
		Str("latest_sync_event_block", fmt.Sprintf("%d", latestLog.BlockNumber)).
		Str("latest_sync_event_tx_hash", latestLog.TxHash.Hex()).
		Msg("Uniswap V2 swap estimation details")

	return &ReservePair{
		Reserve0: reserve0,
		Reserve1: reserve1,
	}, nil
}
