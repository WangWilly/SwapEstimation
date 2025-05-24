package eth

////////////////////////////////////////////////////////////////////////////////

type Config struct {
	BlockRangeSize uint64 `env:"BLOCK_RANGE_SIZE,default=9900"`
}

type client struct {
	cfg Config

	gethClient GethClient
}

func New(cfg Config, gethClient GethClient) *client {
	return &client{
		cfg:        cfg,
		gethClient: gethClient,
	}
}
