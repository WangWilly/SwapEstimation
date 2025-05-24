package ctrlutils

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsValidAddr(t *testing.T) {
	Convey("Given the Ethereum address validation function", t, func() {
		Convey("When validating correctly formatted addresses", func() {
			Convey("With a lowercase address", func() {
				result := IsValidAddr("0x1234567890abcdef1234567890abcdef12345678")
				So(result, ShouldBeTrue)
			})

			Convey("With an uppercase address", func() {
				result := IsValidAddr("0x1234567890ABCDEF1234567890ABCDEF12345678")
				So(result, ShouldBeTrue)
			})

			Convey("With a mixed-case address", func() {
				result := IsValidAddr("0x1234567890aBcDeF1234567890aBcDeF12345678")
				So(result, ShouldBeTrue)
			})
		})

		Convey("When validating incorrectly formatted addresses", func() {
			Convey("With an address that is too short", func() {
				result := IsValidAddr("0x1234567890abcdef1234567890abcdef123456")
				So(result, ShouldBeFalse)
			})

			Convey("With an address that is too long", func() {
				result := IsValidAddr("0x1234567890abcdef1234567890abcdef123456789")
				So(result, ShouldBeFalse)
			})

			Convey("With an address missing 0x prefix", func() {
				result := IsValidAddr("1234567890abcdef1234567890abcdef12345678")
				So(result, ShouldBeFalse)
			})

			Convey("With an address containing invalid characters", func() {
				result := IsValidAddr("0x1234567890abcdef1234567890abcdefg2345678")
				So(result, ShouldBeFalse)
			})

			Convey("With an empty string", func() {
				result := IsValidAddr("")
				So(result, ShouldBeFalse)
			})
		})
	})
}

func TestIsValidUniV2PairAddr(t *testing.T) {
	Convey("Given the Uniswap V2 pair address validation function", t, func() {
		// Real world token addresses
		weth := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
		usdc := "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"

		// Compute the expected pair address
		expectedPairAddr := ComputeUniV2PairAddrStr(weth, usdc)

		Convey("When validating a correct pair address", func() {
			result := IsValidUniV2PairAddr(weth, usdc, expectedPairAddr)

			Convey("Then it should return true", func() {
				So(result, ShouldBeTrue)
			})
		})

		Convey("When validating with reversed token order", func() {
			result := IsValidUniV2PairAddr(usdc, weth, expectedPairAddr)

			Convey("Then it should still return true", func() {
				So(result, ShouldBeTrue)
			})
		})

		Convey("When validating with a different case in the pair address", func() {
			upperPairAddr := expectedPairAddr[:2] + strings.ToUpper(expectedPairAddr[2:])
			result := IsValidUniV2PairAddr(weth, usdc, upperPairAddr)

			Convey("Then it should return true due to case-insensitive comparison", func() {
				So(result, ShouldBeTrue)
			})
		})

		Convey("When validating an incorrect pair address", func() {
			invalidPairAddr := "0x0000000000000000000000000000000000000000"
			result := IsValidUniV2PairAddr(weth, usdc, invalidPairAddr)

			Convey("Then it should return false", func() {
				So(result, ShouldBeFalse)
			})
		})
	})
}
