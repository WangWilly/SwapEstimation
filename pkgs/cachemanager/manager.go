package cachemanager

import (
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

////////////////////////////////////////////////////////////////////////////////

type manager struct {
	clientID    string
	redisClient *redis.Client
}

func New(redisClient *redis.Client) *manager {
	clientID := uuid.New().String()

	return &manager{
		clientID:    clientID,
		redisClient: redisClient,
	}
}
