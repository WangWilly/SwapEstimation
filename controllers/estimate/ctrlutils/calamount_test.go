package ctrlutils

import (
	"math/big"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCalOutAmount(t *testing.T) {
	Convey("Given the CalOutAmount function for token swap calculations", t, func() {
		// Common test data
		srcAddr := "0x1000000000000000000000000000000000000000"
		dstAddr := "0x2000000000000000000000000000000000000000"
		reserve0 := big.NewInt(1000000)
		reserve1 := big.NewInt(2000000)

		Convey("When calculating with valid inputs", func() {
			amountIn := big.NewInt(1000)
			result := CalOutAmount(srcAddr, dstAddr, amountIn, reserve0, reserve1)

			Convey("Then the output amount should be correctly calculated", func() {
				// Expected output based on the formula with 0.3% fee
				expected := big.NewInt(1992) // (1000 * 997 * 2000000) / (1000 * 1000000 + 997 * 1000)
				So(result.Cmp(expected), ShouldEqual, 0)
			})
		})

		Convey("When source and destination tokens are the same", func() {
			amountIn := big.NewInt(1000)
			result := CalOutAmount(srcAddr, srcAddr, amountIn, reserve0, reserve1)

			Convey("Then the output should equal the input", func() {
				So(result.Cmp(amountIn), ShouldEqual, 0)
			})
		})

		Convey("When the token addresses are reversed", func() {
			amountIn := big.NewInt(1000)
			result1 := CalOutAmount(srcAddr, dstAddr, amountIn, reserve0, reserve1)
			result2 := CalOutAmount(dstAddr, srcAddr, amountIn, reserve0, reserve1)

			Convey("Then the results should be different", func() {
				So(result1.Cmp(result2), ShouldNotEqual, 0)
			})
		})

		Convey("When input amount is zero", func() {
			amountIn := big.NewInt(0)
			result := CalOutAmount(srcAddr, dstAddr, amountIn, reserve0, reserve1)

			Convey("Then output should be zero", func() {
				So(result.Cmp(big.NewInt(0)), ShouldEqual, 0)
			})
		})

		Convey("When testing edge cases", func() {
			Convey("With empty source address", func() {
				result := CalOutAmount("", dstAddr, big.NewInt(1000), reserve0, reserve1)
				So(result, ShouldBeNil)
			})

			Convey("With empty destination address", func() {
				result := CalOutAmount(srcAddr, "", big.NewInt(1000), reserve0, reserve1)
				So(result, ShouldBeNil)
			})

			Convey("With nil amount in", func() {
				result := CalOutAmount(srcAddr, dstAddr, nil, reserve0, reserve1)
				So(result, ShouldBeNil)
			})

			Convey("With nil reserve values", func() {
				result := CalOutAmount(srcAddr, dstAddr, big.NewInt(1000), nil, reserve1)
				So(result, ShouldBeNil)

				result = CalOutAmount(srcAddr, dstAddr, big.NewInt(1000), reserve0, nil)
				So(result, ShouldBeNil)
			})

			Convey("With negative amount in", func() {
				result := CalOutAmount(srcAddr, dstAddr, big.NewInt(-1000), reserve0, reserve1)
				So(result, ShouldBeNil)
			})

			Convey("With negative reserve values", func() {
				result := CalOutAmount(srcAddr, dstAddr, big.NewInt(1000), big.NewInt(-1000), reserve1)
				So(result, ShouldBeNil)

				result = CalOutAmount(srcAddr, dstAddr, big.NewInt(1000), reserve0, big.NewInt(-1000))
				So(result, ShouldBeNil)
			})
		})
	})
}
