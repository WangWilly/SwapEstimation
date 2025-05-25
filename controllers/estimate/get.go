package estimate

import (
	"context"
	"fmt"
	"math/big"

	"github.com/WangWilly/swap-estimation/controllers/estimate/ctrlutils"
	"github.com/WangWilly/swap-estimation/pkgs/clients/eth"
	"github.com/WangWilly/swap-estimation/pkgs/clients/ethwss"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

////////////////////////////////////////////////////////////////////////////////

func (c *Controller) Get(ctx *gin.Context) {
	logger := log.Ctx(ctx.Request.Context())
	logger.Debug().Msg("Received estimate request")

	poolAddr, ok := ctx.GetQuery("pool")
	if !ok || poolAddr == "" {
		logger.Error().Msg("Pool address is required")
		ctx.JSON(400, gin.H{"error": "pool address is required"})
		return
	}
	ok = ctrlutils.IsValidAddr(poolAddr)
	if !ok {
		logger.Error().Msg("Invalid pool address format")
		ctx.JSON(400, gin.H{"error": "invalid pool address format"})
		return
	}

	srcTokenAddr, ok := ctx.GetQuery("src")
	if !ok || srcTokenAddr == "" {
		logger.Error().Msg("Source token address is required")
		ctx.JSON(400, gin.H{"error": "source token address is required"})
		return
	}
	ok = ctrlutils.IsValidAddr(srcTokenAddr)
	if !ok {
		logger.Error().Msg("Invalid source token address format")
		ctx.JSON(400, gin.H{"error": "invalid source token address format"})
		return
	}

	destTokenAddr, ok := ctx.GetQuery("dst")
	if !ok || destTokenAddr == "" {
		logger.Error().Msg("Destination token address is required")
		ctx.JSON(400, gin.H{"error": "destination token address is required"})
		return
	}
	ok = ctrlutils.IsValidAddr(destTokenAddr)
	if !ok {
		logger.Error().Msg("Invalid destination token address format")
		ctx.JSON(400, gin.H{"error": "invalid destination token address format"})
		return
	}

	ok = ctrlutils.IsValidUniV2PairAddr(srcTokenAddr, destTokenAddr, poolAddr)
	if !ok {
		logger.Error().Msg("Invalid Uniswap V2 pair address")
		ctx.JSON(400, gin.H{"error": "invalid Uniswap V2 pair address"})
		return
	}

	srcAmountStr, ok := ctx.GetQuery("src_amount")
	if !ok || srcAmountStr == "" {
		logger.Error().Msg("Source amount is required")
		ctx.JSON(400, gin.H{"error": "source amount is required"})
		return
	}
	srcAmount := new(big.Int)
	srcAmount, ok = srcAmount.SetString(srcAmountStr, 10)
	if !ok {
		logger.Error().Msg("Invalid source amount")
		ctx.JSON(400, gin.H{"error": "invalid source amount"})
		return
	}

	////////////////////////////////////////////////////////////////////////////

	// Get the reserve pair from cache or fetch it
	reservePair, err := c.getPair(ctx.Request.Context(), poolAddr)
	if err != nil {
		logger.Error().
			Err(err).
			Str("pool_address", poolAddr).
			Msg("Failed to get Uniswap V2 reserve pair")
		ctx.JSON(500, gin.H{"error": "failed to get reserve pair"})
		return
	}

	estimatedAmount := ctrlutils.CalOutAmount(
		srcTokenAddr,
		destTokenAddr,
		srcAmount,
		reservePair.Reserve0,
		reservePair.Reserve1,
	)
	if estimatedAmount == nil {
		logger.Error().Msg("Failed to calculate output amount")
		ctx.JSON(500, gin.H{"error": "failed to calculate output amount"})
		return
	}

	// plain text response
	ctx.Writer.Header().Set("Content-Type", "text/plain")
	ctx.String(200, estimatedAmount.String())
}

////////////////////////////////////////////////////////////////////////////////

func (c *Controller) getPair(ctx context.Context, poolAddr string) (*eth.ReservePair, error) {
	logger := log.Ctx(ctx)
	logger.Debug().
		Str("pool_address", poolAddr).
		Msg("Getting Uniswap V2 pair for reserve updates")

	pair := c.ethWssClient.GetPair(ctx, poolAddr)
	if pair != nil {
		logger.Debug().
			Str("pool_address", poolAddr).
			Msg("Pair found in cache")
		return (*eth.ReservePair)(pair), nil
	}

	// Use singleflight to prevent duplicate requests for the same estimation
	singleflightKey := "estimate_" + poolAddr
	res, err, _ := c.g4GetEstimate.Do(singleflightKey, func() (any, error) {
		return c.ethClient.UniV2ReservePair(ctx, poolAddr)
	})
	if err != nil {
		logger.Error().
			Err(err).
			Str("pool_address", poolAddr).
			Msg("Failed to get Uniswap V2 reserve pair")
		return nil, err
	}
	if res == nil {
		logger.Error().
			Str("pool_address", poolAddr).
			Msg("Uniswap V2 reserve pair not found")
		return nil, err
	}
	logger.Debug().
		Str("pool_address", poolAddr).
		Msg("Uniswap V2 reserve pair found")

	currPair, ok := res.(*eth.ReservePair)
	if !ok {
		logger.Error().
			Str("pool_address", poolAddr).
			Msg("Unexpected response type from UniV2ReservePair")
		return nil, fmt.Errorf("unexpected response type from UniV2ReservePair: %T", res)
	}

	c.ethWssClient.RegPair(context.Background(), poolAddr, (*ethwss.ReservePair)(currPair))
	return currPair, nil
}
