package employeeattendancerepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"gorm.io/gorm"
)

func (r *repo) Last(ctx context.Context, tx *gorm.DB, employeeID int64) (*models.EmployeeAttendance, error) {
	// Create a variable to hold the result
	var employeeAttendance models.EmployeeAttendance

	// Execute the query
	if err := tx.Where("employee_id = ?", employeeID).
		Order("created_at DESC").
		First(&employeeAttendance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get last employee attendance: %w", err)
	}

	// Return the result
	return &employeeAttendance, nil
}
