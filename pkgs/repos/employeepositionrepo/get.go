package employeepositionrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"gorm.io/gorm"
)

////////////////////////////////////////////////////////////////////////////////

func (r *repo) Get(ctx context.Context, tx *gorm.DB, id int64) (*models.EmployeePosition, error) {
	// Create a variable to hold the result
	var employeePosition models.EmployeePosition

	// Execute the query
	if err := tx.Where("id = ?", id).First(&employeePosition).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get employee position: %w", err)
	}

	// Return the result
	return &employeePosition, nil
}

func (r *repo) MustGet(ctx context.Context, tx *gorm.DB, id int64) (*models.EmployeePosition, error) {
	employeePosition, err := r.Get(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee position: %w", err)
	}
	if employeePosition == nil {
		return nil, fmt.Errorf("employee position not found")
	}

	return employeePosition, nil
}

////////////////////////////////////////////////////////////////////////////////

func (r *repo) GetCurrentByEmployeeID(ctx context.Context, tx *gorm.DB, employeeID int64, nowtime time.Time) (*models.EmployeePosition, error) {
	// Create a variable to hold the result
	var employeePosition models.EmployeePosition

	// Execute the query
	if err := tx.Where("employee_id = ? AND start_date <= ?", employeeID, nowtime).
		Order("start_date DESC").
		First(&employeePosition).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get current employee position: %w", err)
	}

	// Return the result
	return &employeePosition, nil
}
