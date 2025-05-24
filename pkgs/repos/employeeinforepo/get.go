package employeeinforepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"gorm.io/gorm"
)

func (r *repo) Get(ctx context.Context, tx *gorm.DB, id int64) (*models.EmployeeInfo, error) {
	// Create a variable to hold the result
	var employeeInfo models.EmployeeInfo

	// Execute the query
	if err := tx.Where("id = ?", id).First(&employeeInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get employee info: %w", err)
	}

	// Return the result
	return &employeeInfo, nil
}

func (r *repo) MustGet(ctx context.Context, tx *gorm.DB, id int64) (*models.EmployeeInfo, error) {
	employeeInfo, err := r.Get(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee info: %w", err)
	}
	if employeeInfo == nil {
		return nil, fmt.Errorf("employee info not found")
	}

	return employeeInfo, nil
}
