package ctrlutils

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

////////////////////////////////////////////////////////////////////////////////

func ComputePairAddr(factory, tokenA, tokenB common.Address, initCodeHash common.Hash) common.Address {
	// Sort token addresses to determine token0 and token1
	var token0, token1 common.Address
	if strings.ToLower(tokenA.Hex()) < strings.ToLower(tokenB.Hex()) {
		token0 = tokenA
		token1 = tokenB
	} else {
		token0 = tokenB
		token1 = tokenA
	}

	// Compute the salt: keccak256(abi.encodePacked(token0, token1))
	salt := crypto.Keccak256Hash(
		append(token0.Bytes(), token1.Bytes()...),
	)

	// Compute the CREATE2 address using the formula:
	// address = keccak256(0xff ++ factory ++ salt ++ init_code_hash)[12:]
	data := []byte{0xff}
	data = append(data, factory.Bytes()...)
	data = append(data, salt.Bytes()...)
	data = append(data, initCodeHash.Bytes()...)
	hash := crypto.Keccak256Hash(data)

	// The last 20 bytes of the hash represent the address
	return common.BytesToAddress(hash.Bytes()[12:])
}

func ComputePairAddrStr(factoryStr, tokenAStr, tokenBStr, initCodeHashStr string) string {
	factory := common.HexToAddress(factoryStr)
	tokenA := common.HexToAddress(tokenAStr)
	tokenB := common.HexToAddress(tokenBStr)
	initCodeHash := common.HexToHash(initCodeHashStr)

	pairAddr := ComputePairAddr(factory, tokenA, tokenB, initCodeHash)
	return pairAddr.Hex()
}

func ComputeUniV2PairAddrStr(tokenAStr, tokenBStr string) string {
	// Uniswap V2 factory address on Ethereum mainnet
	factory := common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f")

	// Uniswap V2 pair contract init code hash
	initCodeHash := common.HexToHash("0x96e8ac4277198ff8b6f785478aa9a39f403cb768dd02cbee326c3e7da348845f")

	return ComputePairAddrStr(factory.Hex(), tokenAStr, tokenBStr, initCodeHash.Hex())
}
