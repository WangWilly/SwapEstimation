package ctrlutils

import "strings"

////////////////////////////////////////////////////////////////////////////////

func IsValidAddr(addr string) bool {
	if len(addr) != 42 {
		return false
	}
	if addr[:2] != "0x" {
		return false
	}
	for _, c := range addr[2:] {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func IsValidUniV2PairAddr(tokenAStr, tokenBStr, pairAddrStr string) bool {
	computedAddr := ComputeUniV2PairAddrStr(tokenAStr, tokenBStr)
	return strings.EqualFold(computedAddr, pairAddrStr)
}
