package cachemanager

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

////////////////////////////////////////////////////////////////////////////////

type cacheMainKey string

const (
	employeeDetailV1 cacheMainKey = "employee_detail_v1"
	attendanceV1     cacheMainKey = "attendance_v1"
)

func buildCacheFullKey(k cacheMainKey, pairs map[string]any) (string, error) {
	if len(pairs) == 0 {
		return string(k), nil
	}

	intergatedPairs := lo.Filter(lo.MapToSlice(pairs, func(key string, value any) string {
		switch value := value.(type) {
		case int:
			return fmt.Sprintf("%s-%d", key, value)
		case int64:
			return fmt.Sprintf("%s-%d", key, value)
		case string:
			return fmt.Sprintf("%s-%s", key, value)
		case bool:
			return fmt.Sprintf("%s-%t", key, value)
		default:
			return ""
		}
	}), func(pair string, _ int) bool {
		return pair != ""
	})

	if len(intergatedPairs) == 0 {
		return "", fmt.Errorf("no valid pairs found")
	}

	return "[" + string(k) + "]" + strings.Join(intergatedPairs, ":"), nil

}

////////////////////////////////////////////////////////////////////////////////

func gobEncode(value any) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func gobDecode(data []byte, value any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(value)
}

////////////////////////////////////////////////////////////////////////////////

func setItem[K any](redisClient *redis.Client, ctx context.Context, key string, item K, expired time.Duration) error {

	serializedItems, err := gobEncode(item)
	if err != nil {
		return err
	}

	if err = redisClient.Set(ctx, key, serializedItems, expired).Err(); err != nil {
		return err
	}

	return nil
}

func getItem[K any](redisClient *redis.Client, ctx context.Context, key string, must bool) (*K, error) {

	serializedItem, err := redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			if must {
				return nil, fmt.Errorf("cache miss: %s", key)
			}

			return nil, nil
		}
		return nil, err
	}

	item := new(K)
	err = gobDecode(serializedItem, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}
