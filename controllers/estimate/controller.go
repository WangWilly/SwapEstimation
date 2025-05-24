package estimate

import (
	"github.com/gin-gonic/gin"

	"golang.org/x/sync/singleflight"
)

////////////////////////////////////////////////////////////////////////////////

type Config struct {
}

type Controller struct {
	cfg Config

	ethClient EthClient

	g4GetEstimate *singleflight.Group
}

func NewController(
	cfg Config,
	ethClient EthClient,
) *Controller {
	g4GetEstimate := &singleflight.Group{}

	return &Controller{
		cfg:           cfg,
		ethClient:     ethClient,
		g4GetEstimate: g4GetEstimate,
	}
}

func (c *Controller) RegisterRoutes(r *gin.Engine) {
	////////////////////////////////////////////////////////////////////////////
	// price estimation
	r.GET("/estimate", c.Get)
}
