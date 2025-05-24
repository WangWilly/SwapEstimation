package cachemanager

import (
	"context"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/dtos"
)

////////////////////////////////////////////////////////////////////////////////

func (m *manager) SetEmployeeDetailV1(
	ctx context.Context,
	employeeID int64,
	data dtos.EmployeeV1Response,
	expired time.Duration,
) error {
	fullKey, err := buildCacheFullKey(employeeDetailV1, map[string]any{
		"employee_id": employeeID,
	})
	if err != nil {
		return err
	}

	return setItem(m.redisClient, ctx, fullKey, data, expired)
}

func (m *manager) GetEmployeeDetailV1(
	ctx context.Context,
	employeeID int64,
) (*dtos.EmployeeV1Response, error) {
	fullKey, err := buildCacheFullKey(employeeDetailV1, map[string]any{
		"employee_id": employeeID,
	})
	if err != nil {
		return nil, err
	}

	return getItem[dtos.EmployeeV1Response](m.redisClient, ctx, fullKey, false)
}

func (m *manager) DeleteEmployeeDetailV1(
	ctx context.Context,
	employeeID int64,
) error {
	fullKey, err := buildCacheFullKey(employeeDetailV1, map[string]any{
		"employee_id": employeeID,
	})
	if err != nil {
		return err
	}

	return m.redisClient.Del(ctx, fullKey).Err()
}
