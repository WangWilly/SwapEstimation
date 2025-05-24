package cachemanager

import (
	"context"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/dtos"
)

////////////////////////////////////////////////////////////////////////////////

func (m *manager) SetAttendanceV1(
	ctx context.Context,
	employeeID int64,
	data dtos.AttendanceV1Response,
	expired time.Duration,
) error {
	fullKey, err := buildCacheFullKey(attendanceV1, map[string]any{
		"employee_id": employeeID,
	})
	if err != nil {
		return err
	}

	return setItem(m.redisClient, ctx, fullKey, data, expired)
}

func (m *manager) GetAttendanceV1(
	ctx context.Context,
	employeeID int64,
) (*dtos.AttendanceV1Response, error) {
	fullKey, err := buildCacheFullKey(attendanceV1, map[string]any{
		"employee_id": employeeID,
	})
	if err != nil {
		return nil, err
	}

	return getItem[dtos.AttendanceV1Response](m.redisClient, ctx, fullKey, false)
}

func (m *manager) DeleteAttendanceV1(
	ctx context.Context,
	employeeID int64,
) error {
	fullKey, err := buildCacheFullKey(attendanceV1, map[string]any{
		"employee_id": employeeID,
	})
	if err != nil {
		return err
	}

	return m.redisClient.Del(ctx, fullKey).Err()
}
