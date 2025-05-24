package ctrlutils

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/smartystreets/goconvey/convey"
)

func TestComputePairAddr(t *testing.T) {
	Convey("Given the Uniswap V2 pair address computation function", t, func() {
		// Real Ethereum mainnet addresses
		factory := common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f")
		weth := common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
		usdc := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
		initCodeHash := common.HexToHash("0x96e8ac4277198ff8b6f785478aa9a39f403cb768dd02cbee326c3e7da348845f")

		// Known WETH-USDC pair address on mainnet
		expectedPairAddr := common.HexToAddress("0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc")

		Convey("When computing pair address with token order A->B", func() {
			result := ComputePairAddr(factory, weth, usdc, initCodeHash)

			Convey("Then the computed address should match the expected pair address", func() {
				So(result.Hex(), ShouldEqual, expectedPairAddr.Hex())
			})
		})

		Convey("When computing pair address with token order B->A", func() {
			result := ComputePairAddr(factory, usdc, weth, initCodeHash)

			Convey("Then the computed address should still match the expected pair address", func() {
				So(result.Hex(), ShouldEqual, expectedPairAddr.Hex())
			})
		})

		Convey("When comparing results with different token order", func() {
			result1 := ComputePairAddr(factory, weth, usdc, initCodeHash)
			result2 := ComputePairAddr(factory, usdc, weth, initCodeHash)

			Convey("Then both results should be identical", func() {
				So(result1.Hex(), ShouldEqual, result2.Hex())
			})
		})
	})
}

func TestComputePairAddrStr(t *testing.T) {
	Convey("Given the string-based pair address computation function", t, func() {
		factoryStr := "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"
		wethStr := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
		usdcStr := "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
		initCodeHashStr := "0x96e8ac4277198ff8b6f785478aa9a39f403cb768dd02cbee326c3e7da348845f"

		// Known WETH-USDC pair address on mainnet
		expectedPairAddrStr := "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc"

		Convey("When computing pair address with string inputs", func() {
			result := ComputePairAddrStr(factoryStr, wethStr, usdcStr, initCodeHashStr)

			Convey("Then the computed address should match the expected pair address", func() {
				So(strings.EqualFold(result, expectedPairAddrStr), ShouldBeTrue)
			})
		})

		Convey("When computing with reversed token order", func() {
			result := ComputePairAddrStr(factoryStr, usdcStr, wethStr, initCodeHashStr)

			Convey("Then the computed address should still match the expected pair address", func() {
				So(strings.EqualFold(result, expectedPairAddrStr), ShouldBeTrue)
			})
		})

		Convey("When comparing results for both token orders", func() {
			result1 := ComputePairAddrStr(factoryStr, wethStr, usdcStr, initCodeHashStr)
			result2 := ComputePairAddrStr(factoryStr, usdcStr, wethStr, initCodeHashStr)

			Convey("Then both results should be identical", func() {
				So(strings.EqualFold(result1, result2), ShouldBeTrue)
			})
		})
	})
}

func TestComputeUniV2PairAddrStr(t *testing.T) {
	Convey("Given the Uniswap V2 specific pair computation function", t, func() {
		wethStr := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
		usdcStr := "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"

		// Known WETH-USDC pair address on Uniswap V2
		expectedPairAddrStr := "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc"

		Convey("When computing Uniswap V2 pair address", func() {
			result := ComputeUniV2PairAddrStr(wethStr, usdcStr)

			Convey("Then the computed address should match the expected Uniswap V2 pair address", func() {
				So(strings.EqualFold(result, expectedPairAddrStr), ShouldBeTrue)
			})
		})

		Convey("When computing with reversed token order", func() {
			result := ComputeUniV2PairAddrStr(usdcStr, wethStr)

			Convey("Then the computed address should still be correct", func() {
				So(strings.EqualFold(result, expectedPairAddrStr), ShouldBeTrue)
			})
		})

		Convey("When comparing results for both orders", func() {
			result1 := ComputeUniV2PairAddrStr(wethStr, usdcStr)
			result2 := ComputeUniV2PairAddrStr(usdcStr, wethStr)

			Convey("Then both results should be identical", func() {
				So(strings.EqualFold(result1, result2), ShouldBeTrue)
			})
		})
	})
}
