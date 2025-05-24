package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WangWilly/swap-estimation/controllers/estimate"
	"github.com/WangWilly/swap-estimation/pkgs/clients/eth"
	"github.com/WangWilly/swap-estimation/pkgs/middleware"
	"github.com/WangWilly/swap-estimation/pkgs/utils"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/sethvargo/go-envconfig"
)

////////////////////////////////////////////////////////////////////////////////

type envConfig struct {
	// Server configuration
	Port string `env:"PORT,default=8080"`
	Host string `env:"HOST,default=0.0.0.0"`

	// Eth client configuration
	GethClientURL string `env:"GETH_CLIENT_URL,required"`

	EthClientCfg eth.Config `env:",prefix=ETH_CLIENT_"`
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	ctx := context.Background()
	utils.InitLogging(ctx)
}

func main() {
	ctx := context.Background()
	logger := utils.GetDetailedLogger().With().Caller().Logger()

	// Load environment variables
	cfg := &envConfig{}
	err := envconfig.Process(ctx, cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load environment variables")
	}

	////////////////////////////////////////////////////////////////////////////
	// Initialize Gin router

	r := utils.GetDefaultRouter()
	r.Use(middleware.LoggingMiddleware())

	////////////////////////////////////////////////////////////////////////////
	// Initialize modules

	gethClient, err := ethclient.Dial(cfg.GethClientURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	ethClient := eth.New(cfg.EthClientCfg, gethClient)

	////////////////////////////////////////////////////////////////////////////
	// Initialize the controllers

	estimateCtrlCfg := estimate.Config{}
	estimateCtrl := estimate.NewController(
		estimateCtrlCfg,
		ethClient,
	)
	estimateCtrl.RegisterRoutes(r)

	////////////////////////////////////////////////////////////////////////////

	// Set up the server
	srv := &http.Server{
		Addr:    cfg.Host + ":" + cfg.Port,
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	////////////////////////////////////////////////////////////////////////////

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	// Kill (no param) default sends syscall.SIGTERM
	// Kill -2 is syscall.SIGINT
	// Kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Log shutdown message
	logger.Info().Msg("Received shutdown signal, shutting down server...")
	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	logger.Info().Msg("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		// Handle shutdown error
		logger.Fatal().Err(err).Msg("Failed to shutdown server")
	}

	// Wait for tasks to finish or timeout
	<-ctx.Done()
	logger.Info().Msg("Server shutdown complete.")
}
