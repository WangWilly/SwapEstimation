package testutils

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/WangWilly/swap-estimation/pkgs/utils"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/redis/go-redis/v9"
	"github.com/sethvargo/go-envconfig"
)

////////////////////////////////////////////////////////////////////////////////

type DockerTestRedisInstance struct {
	DockerPool      *dockertest.Pool
	DockerContainer *dockertest.Resource

	RedisClient *redis.Client
}

var redisInstance *DockerTestRedisInstance

////////////////////////////////////////////////////////////////////////////////

func (r *DockerTestRedisInstance) MustClose() {
	if err := r.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func (r *DockerTestRedisInstance) Close() (retErr error) {
	defer func() {
		if err := r.DockerPool.Purge(r.DockerContainer); err != nil {
			retErr = fmt.Errorf("failed to purge database container: %w", err)
			return
		}
	}()

	if err := r.RedisClient.Close(); err != nil {
		return fmt.Errorf("failed to close redis client: %w", err)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func NewRedis() (*DockerTestRedisInstance, error) {
	////////////////////////////////////////////////////////////////////////////
	// Load environment variables

	ctx := context.Background()
	redisCfg := utils.RedisConfig{}
	if err := envconfig.Process(ctx, &redisCfg); err != nil {
		return nil, fmt.Errorf("failed to process env config: %w", err)
	}

	dockerTestOptions := &dockertest.RunOptions{
		Repository: "redis",
		Tag:        "7.4.3-alpine",
		Env: []string{
			"LANG=C",
		},
	}
	dockerTestPort := "6379/tcp"

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
	_ = container.Expire(120)

	////////////////////////////////////////////////////////////////////////////
	// Wait for the container to be ready

	addr := container.GetBoundIP(dockerTestPort) +
		":" + container.GetPort(dockerTestPort)
	redisCfg.Addr = addr

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	var redisClient *redis.Client
	if err = pool.Retry(func() error {
		redisClient, err = utils.GetRedis(ctx, redisCfg)
		if err != nil {
			return fmt.Errorf("failed to create redis connect => %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &DockerTestRedisInstance{
		DockerPool:      pool,
		DockerContainer: container,
		RedisClient:     redisClient,
	}, nil
}

func MustNewRedis() *DockerTestRedisInstance {
	instance, err := NewRedis()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	return instance
}

func GetRedis() *DockerTestRedisInstance {
	return redisInstance
}

////////////////////////////////////////////////////////////////////////////////

func BeforeTestDbRedis(m *testing.M) {
	redisInstance = MustNewRedis()

	code := m.Run()
	// defer() won't work since os.Exit() is called
	if err := redisInstance.Close(); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
