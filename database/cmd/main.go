package main

import (
	"context"

	"github.com/WangWilly/labs-hr-go/pkgs/utils"
	"github.com/sethvargo/go-envconfig"

	"github.com/WangWilly/labs-hr-go/database/migrations"
)

////////////////////////////////////////////////////////////////////////////////

type envConfig struct {
	DbCfg utils.DbConfig `env:",prefix="`
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	ctx := context.Background()

	// Load environment variables
	cfg := &envConfig{}
	err := envconfig.Process(ctx, cfg)
	if err != nil {
		panic(err)
	}

	////////////////////////////////////////////////////////////////////////////
	// setup database
	db, err := utils.GetDB(cfg.DbCfg)
	if err != nil {
		panic(err)
	}

	////////////////////////////////////////////////////////////////////////////
	// run migrations

	if err := migrations.Apply(db); err != nil {
		panic(err)
	}
}
