package estimate

import (
	"context"

	"github.com/WangWilly/swap-estimation/pkgs/clients/eth"
	"github.com/WangWilly/swap-estimation/pkgs/clients/ethwss"
)

//go:generate mockgen -source=interface.go -destination=interface_mock.go -package=estimate
type EthClient interface {
	UniV2ReservePair(ctx context.Context, pairAddrStr string) (*eth.ReservePair, error)
}

type EthWssClient interface {
	GetPair(ctx context.Context, address string) *ethwss.ReservePair
	RegPair(ctx context.Context, address string, initPair *ethwss.ReservePair) error
}
