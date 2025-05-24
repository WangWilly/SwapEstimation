package attendance

import (
	"context"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/dtos"
	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"gorm.io/gorm"
)

//go:generate mockgen -source=interface.go -destination=interface_mock.go -package=attendance
type TimeModule interface {
	Now() time.Time
}

type EmployeePositionRepo interface {
	GetCurrentByEmployeeID(ctx context.Context, tx *gorm.DB, employeeID int64, nowtime time.Time) (*models.EmployeePosition, error)
}

type EmployeeAttendanceRepo interface {
	CreateForClockIn(ctx context.Context, tx *gorm.DB, employeeID int64, positionID int64, clockInTime time.Time) (*models.EmployeeAttendance, error)
	Last(ctx context.Context, tx *gorm.DB, employeeID int64) (*models.EmployeeAttendance, error)
	UpdateForClockOut(ctx context.Context, tx *gorm.DB, attendanceID int64, clockOutTime time.Time) (*models.EmployeeAttendance, error)
}

type CacheManager interface {
	GetAttendanceV1(ctx context.Context, employeeID int64) (*dtos.AttendanceV1Response, error)
	GetEmployeeDetailV1(ctx context.Context, employeeID int64) (*dtos.EmployeeV1Response, error)
	SetAttendanceV1(ctx context.Context, employeeID int64, data dtos.AttendanceV1Response, expired time.Duration) error
}
