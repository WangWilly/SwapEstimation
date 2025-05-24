package testutils

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/WangWilly/labs-hr-go/database/migrations"
	"github.com/WangWilly/labs-hr-go/pkgs/utils"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/sethvargo/go-envconfig"
)

////////////////////////////////////////////////////////////////////////////////

type DockerTestDatabaseInstance struct {
	DockerPool      *dockertest.Pool
	DockerContainer *dockertest.Resource

	DB *gorm.DB
}

var dbInstance *DockerTestDatabaseInstance

////////////////////////////////////////////////////////////////////////////////

func (d *DockerTestDatabaseInstance) MustClose() {
	if err := d.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func (d *DockerTestDatabaseInstance) Close() (retErr error) {
	defer func() {
		if err := d.DockerPool.Purge(d.DockerContainer); err != nil {
			retErr = fmt.Errorf("failed to purge database container: %w", err)
			return
		}
	}()

	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sqlDB from gorm: %w", err)
	}

	if err = sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close sqlDB: %w", err)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func NewDatabase() (*DockerTestDatabaseInstance, error) {
	////////////////////////////////////////////////////////////////////////////
	// Load environment variables

	ctx := context.Background()

	dbCfg := &utils.DbConfig{}
	if err := envconfig.Process(ctx, dbCfg); err != nil {
		return nil, fmt.Errorf("failed to process env config: %w", err)
	}

	var dockerTestOptions *dockertest.RunOptions
	var dockerTestPort string
	switch dbCfg.Driver {
	case "mysql":
		dockerTestOptions = &dockertest.RunOptions{
			Repository: "mysql",
			Tag:        "9.3.0",
			Env: []string{
				"MYSQL_ROOT_PASSWORD=" + dbCfg.Password,
				"MYSQL_DATABASE=" + dbCfg.Database,
				"MYSQL_USER=" + dbCfg.User,
				"MYSQL_PASSWORD=" + dbCfg.Password,
			},
		}
		dockerTestPort = "3306/tcp"
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", dbCfg.Driver)
	}

	////////////////////////////////////////////////////////////////////////////
	// Create a new pool

	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("failed to create docker pool: %w", err)
	}

	container, err := pool.RunWithOptions(dockerTestOptions, func(c *docker.HostConfig) {
		c.AutoRemove = true
		c.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start database container: %w", err)
	}
	_ = container.Expire(1200)

	////////////////////////////////////////////////////////////////////////////
	// Wait for the container to be ready

	host := container.GetBoundIP(dockerTestPort)
	port := container.GetPort(dockerTestPort)
	dbCfg.Host = host
	dbCfg.Port = port

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	var gormDB *gorm.DB
	if err = pool.Retry(func() error {
		gormDB, err = utils.GetDB(*dbCfg)
		if err != nil {
			return fmt.Errorf("failed to create database connect: %w", err)
		}

		sqlDB, err := gormDB.DB()
		if err != nil {
			return fmt.Errorf("failed to get sql db: %w", err)
		}

		if err := migrations.Apply(gormDB); err != nil {
			return fmt.Errorf("failed to create migrate driver: %w", err)
		}

		return sqlDB.Ping()
	}); err != nil {
		return nil, err
	}

	return &DockerTestDatabaseInstance{
		DockerPool:      pool,
		DockerContainer: container,
		DB:              gormDB,
	}, nil
}

func MustNewDatabase() *DockerTestDatabaseInstance {
	instance, err := NewDatabase()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	return instance
}

func GetDB() *DockerTestDatabaseInstance {
	return dbInstance
}

////////////////////////////////////////////////////////////////////////////////

func BeforeTestDb(m *testing.M) {
	dbInstance = MustNewDatabase()

	code := m.Run()
	// defer() won't work since os.Exit() is called
	if err := dbInstance.Close(); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
