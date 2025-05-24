package estimate

import (
	"context"

	"github.com/WangWilly/swap-estimation/pkgs/clients/eth"
)

//go:generate mockgen -source=interface.go -destination=interface_mock.go -package=estimate
type EthClient interface {
	UniV2ReservePair(ctx context.Context, pairAddrStr string) (*eth.ReservePair, error)
}
