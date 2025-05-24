package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WangWilly/labs-hr-go/controllers/attendance"
	"github.com/WangWilly/labs-hr-go/controllers/employee"
	"github.com/WangWilly/labs-hr-go/database/migrations"
	"github.com/WangWilly/labs-hr-go/pkgs/cachemanager"
	"github.com/WangWilly/labs-hr-go/pkgs/middleware"
	"github.com/WangWilly/labs-hr-go/pkgs/repos/employeeattendancerepo"
	"github.com/WangWilly/labs-hr-go/pkgs/repos/employeeinforepo"
	"github.com/WangWilly/labs-hr-go/pkgs/repos/employeepositionrepo"
	"github.com/WangWilly/labs-hr-go/pkgs/seed"
	"github.com/WangWilly/labs-hr-go/pkgs/timemodule"
	"github.com/WangWilly/labs-hr-go/pkgs/utils"

	"github.com/sethvargo/go-envconfig"
)

////////////////////////////////////////////////////////////////////////////////

type envConfig struct {
	// Server configuration
	Port string `env:"PORT,default=8080"`
	Host string `env:"HOST,default=0.0.0.0"`

	// Database configuration
	DbCfg     utils.DbConfig `env:",prefix="`
	DbMigrate bool           `env:"DB_MIGRATE,default=true"`
	DbSeed    bool           `env:"DB_SEED,default=false"`

	// Redis configuration
	RedisCfg utils.RedisConfig `env:",prefix="`
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
	// Setup database

	db, err := utils.GetDB(cfg.DbCfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get sqlDB from db")
	}

	if cfg.DbMigrate {
		if err := migrations.Apply(db); err != nil {
			logger.Fatal().Err(err).Msg("Failed to apply database migrations")
		}
		logger.Info().Msg("Database migrations applied successfully")
	}

	// Seed the database if DB_SEED is true
	if cfg.DbSeed {
		logger.Info().Msg("Seeding database with dummy data...")
		if err := seed.SeedData(ctx, db); err != nil {
			logger.Fatal().Err(err).Msg("Failed to seed database")
		}
		logger.Info().Msg("Database seeded successfully")
	}

	////////////////////////////////////////////////////////////////////////////
	// Initialize Redis client

	redisClient, err := utils.GetRedis(ctx, cfg.RedisCfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	logger.Info().Msg("Redis client created successfully!")

	////////////////////////////////////////////////////////////////////////////
	// Initialize modules

	timeModule := timemodule.New()
	employeeInfoRepo := employeeinforepo.New()
	employeePositionRepo := employeepositionrepo.New()
	employeeAttendanceRepo := employeeattendancerepo.New()
	cacheManager := cachemanager.New(redisClient)

	////////////////////////////////////////////////////////////////////////////
	// Initialize the controllers

	employeeCtrlCfg := employee.Config{}
	employeeCtrl := employee.NewController(
		employeeCtrlCfg,
		db,
		timeModule,
		employeeInfoRepo,
		employeePositionRepo,
		cacheManager,
	)
	employeeCtrl.RegisterRoutes(r)

	attendanceCtrlCfg := attendance.Config{}
	attendanceCtrl := attendance.NewController(
		attendanceCtrlCfg,
		db,
		timeModule,
		employeePositionRepo,
		employeeAttendanceRepo,
		cacheManager,
	)
	attendanceCtrl.RegisterRoutes(r)

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

	// Shutdown the server gracefully
	if err := redisClient.Close(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to close Redis client")
	}
	logger.Info().Msg("Redis client closed successfully!")

	// Close the database connection
	if err := sqlDB.Close(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to close database connection")
	}
	logger.Info().Msg("Database connection closed successfully!")

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
