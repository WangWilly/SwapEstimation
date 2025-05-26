package ethwss

import (
	"context"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

////////////////////////////////////////////////////////////////////////////////

const uniswapV2PairABI = `[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint112","name":"reserve0","type":"uint112"},{"indexed":false,"internalType":"uint112","name":"reserve1","type":"uint112"}],"name":"Sync","type":"event"}]`

////////////////////////////////////////////////////////////////////////////////

func (c *client) RegPair(ctx context.Context, address string, initPair *ReservePair) error {
	logger := log.Ctx(ctx)

	// Check if this address is already being registered by another thread
	c.addressLock.Lock()
	if c.registeringPairs[address] {
		c.addressLock.Unlock()
		logger.Debug().
			Str("pair_address", address).
			Msg("Registration for this pair already in progress, skipping")
		return nil
	}

	// Mark this address as being registered
	c.registeringPairs[address] = true
	c.addressLock.Unlock()

	// Ensure we remove the flag when we're done
	defer func() {
		c.addressLock.Lock()
		delete(c.registeringPairs, address)
		c.addressLock.Unlock()
	}()

	////////////////////////////////////////////////////////////////////////////

	logger.Debug().
		Str("pair_address", address).
		Msg("Registering Uniswap V2 pair for reserve updates")
	if _, ok := c.reservePairCacheMap[address]; ok {
		logger.Warn().
			Str("pair_address", address).
			Msg("Pair already registered, skipping registration")
		return nil // Pair already registered
	}
	// from histrical event logs
	c.reservePairCacheMap[address] = initPair

	parsedABI, err := abi.JSON(strings.NewReader(uniswapV2PairABI))
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to parse Uniswap V2 Pair ABI")
		return err
	}
	pairAddress := common.HexToAddress(address)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{pairAddress},
		Topics:    [][]common.Hash{{parsedABI.Events["Sync"].ID}},
	}

	logs := make(chan types.Log)
	sub, err := c.gethWssClient.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to subscribe to Uniswap V2 Pair logs")
		return err
	}

	// Save subscription for later management
	c.pairSubscriptions[address] = sub

	// Create and start a timer for this subscription
	timer := time.NewTimer(c.cfg.ListenPairPeriod)
	c.pairTimers[address] = timer

	go func() {
		defer func() {
			// Cleanup when done
			delete(c.reservePairCacheMap, address)
			delete(c.pairSubscriptions, address)
			delete(c.pairTimers, address)
		}()

		for {
			select {
			case <-timer.C:
				logger.Info().
					Str("pair_address", address).
					Msg("Subscription period expired, unsubscribing")
				sub.Unsubscribe()
				return
			case err := <-sub.Err():
				logger.Error().
					Err(err).
					Msg("Subscription error")
				return
			case vLog := <-logs:
				var event ReservePair
				err := parsedABI.UnpackIntoInterface(&event, "Sync", vLog.Data)
				if err != nil {
					log.Printf("Failed to unpack log: %v", err)
					continue
				}
				c.reservePairCacheMap[address] = &event
			case <-ctx.Done():
				logger.Info().
					Msg("Context done, stopping subscription")
				return
			}
		}
	}()

	return nil
}
