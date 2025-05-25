package ethwss

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

//go:generate mockgen -source=interface.go -destination=interface_mock.go -package=ethwss
type GethWssClient interface {
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
	Close()
}
