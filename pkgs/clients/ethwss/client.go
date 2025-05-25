package ethwss

import (
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/event"
)

////////////////////////////////////////////////////////////////////////////////

type Config struct {
	ListenPairPeriod time.Duration `env:"LISTEN_PAIR_PERIOD,default=2m"`
}

////////////////////////////////////////////////////////////////////////////////

type ReservePair struct {
	Reserve0 *big.Int
	Reserve1 *big.Int
}

////////////////////////////////////////////////////////////////////////////////

type client struct {
	cfg Config

	reservePairCacheMap map[string]*ReservePair
	gethWssClient       GethWssClient

	// Track subscriptions and timers
	pairSubscriptions map[string]event.Subscription
	pairTimers        map[string]*time.Timer

	// Track addresses being registered to prevent concurrent registration of the same address
	addressLock      sync.Mutex
	registeringPairs map[string]bool
}

func New(cfg Config, gethWssClient GethWssClient) *client {
	return &client{
		cfg:                 cfg,
		reservePairCacheMap: make(map[string]*ReservePair),
		gethWssClient:       gethWssClient,
		pairSubscriptions:   make(map[string]event.Subscription),
		pairTimers:          make(map[string]*time.Timer),
		registeringPairs:    make(map[string]bool),
	}
}
