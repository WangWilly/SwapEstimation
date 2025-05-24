package employee

import (
	"context"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/dtos"
	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"gorm.io/gorm"
)

//go:generate mockgen -source=interface.go -destination=interface_mock.go -package=employee
type TimeModule interface {
	Now() time.Time
}

type EmployeeInfoRepo interface {
	Create(ctx context.Context, tx *gorm.DB, data *models.EmployeeInfo) error
	MustGet(ctx context.Context, tx *gorm.DB, id int64) (*models.EmployeeInfo, error)
	Save(ctx context.Context, tx *gorm.DB, data *models.EmployeeInfo) error
}

type EmployeePositionRepo interface {
	Create(ctx context.Context, tx *gorm.DB, data *models.EmployeePosition, nowtime time.Time) error
	Get(ctx context.Context, tx *gorm.DB, id int64) (*models.EmployeePosition, error)
	GetCurrentByEmployeeID(ctx context.Context, tx *gorm.DB, employeeID int64, nowtime time.Time) (*models.EmployeePosition, error)
	MustGet(ctx context.Context, tx *gorm.DB, id int64) (*models.EmployeePosition, error)
}

type CacheManager interface {
	GetEmployeeDetailV1(ctx context.Context, employeeID int64) (*dtos.EmployeeV1Response, error)
	SetEmployeeDetailV1(ctx context.Context, employeeID int64, data dtos.EmployeeV1Response, expired time.Duration) error
	DeleteEmployeeDetailV1(ctx context.Context, employeeID int64) error
}
