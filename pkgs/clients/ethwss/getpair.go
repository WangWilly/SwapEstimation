package ethwss

import (
	"context"

	"github.com/rs/zerolog/log"
)

////////////////////////////////////////////////////////////////////////////////

func (c *client) GetPair(ctx context.Context, address string) *ReservePair {
	logger := log.Ctx(ctx)
	logger.Debug().
		Str("pair_address", address).
		Msg("Getting Uniswap V2 pair for reserve updates")

	if pair, ok := c.reservePairCacheMap[address]; ok {
		logger.Debug().
			Str("pair_address", address).
			Msg("Pair found in cache")

		// Extend the subscription period by resetting the timer
		if timer, exists := c.pairTimers[address]; exists {
			timer.Reset(c.cfg.ListenPairPeriod)
			logger.Debug().
				Str("pair_address", address).
				Dur("period", c.cfg.ListenPairPeriod).
				Msg("Extended subscription period")
		}

		return pair
	}

	logger.Warn().
		Str("pair_address", address).
		Msg("Pair not found in cache, returning nil")
	return nil
}
