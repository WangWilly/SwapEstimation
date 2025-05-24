package employeepositionrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"gorm.io/gorm"
)

////////////////////////////////////////////////////////////////////////////////

func (r *repo) Create(ctx context.Context, tx *gorm.DB, data *models.EmployeePosition, nowtime time.Time) error {
	// The start_date must be greater than the current position's start_date
	currentPosition, err := r.GetCurrentByEmployeeID(ctx, tx, data.EmployeeID, nowtime)
	if err != nil {
		return fmt.Errorf("failed to get current employee position: %w", err)
	}
	if currentPosition != nil && currentPosition.StartDate.After(data.StartDate) {
		return fmt.Errorf("start_date must be greater than the current position's start_date")
	}

	// Create the new employee position
	if err := tx.
		Create(data).Error; err != nil {
		return fmt.Errorf("failed to create employee position: %w", err)
	}

	return nil
}
