package ctrlutils

import (
	"math/big"
)

////////////////////////////////////////////////////////////////////////////////

func CalOutAmount(
	srcTokenAddrStr string,
	dstTokenAddrStr string,
	amountIn, reserve0, reserve1 *big.Int,
) *big.Int {
	if srcTokenAddrStr == "" || dstTokenAddrStr == "" {
		return nil
	}
	if amountIn == nil || reserve0 == nil || reserve1 == nil {
		return nil
	}
	if srcTokenAddrStr == dstTokenAddrStr {
		// If source and destination tokens are the same, return the input amount
		return new(big.Int).Set(amountIn)
	}
	if amountIn.Sign() < 0 || reserve0.Sign() < 0 || reserve1.Sign() < 0 {
		return nil
	}
	if amountIn.Sign() == 0 {
		// If input amount is zero, return zero
		return big.NewInt(0)
	}

	var reserveIn, reserveOut *big.Int
	if srcTokenAddrStr < dstTokenAddrStr {
		reserveIn = reserve0
		reserveOut = reserve1
	} else {
		reserveIn = reserve1
		reserveOut = reserve0
	}

	// Calculate amount out using Uniswap V2 formula
	amountInWithFee := new(big.Int).Mul(amountIn, big.NewInt(997))
	numerator := new(big.Int).Mul(amountInWithFee, reserveOut)
	denominator := new(big.Int).Add(new(big.Int).Mul(reserveIn, big.NewInt(1000)), amountInWithFee)
	amountOut := new(big.Int).Div(numerator, denominator)

	if amountOut.Sign() < 0 {
		// If the calculated amount out is negative, return nil
		return nil
	}

	return amountOut
}
