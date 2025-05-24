package employeeinforepo

import (
	"context"
	"fmt"

	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"gorm.io/gorm"
)

func (r *repo) Save(ctx context.Context, tx *gorm.DB, data *models.EmployeeInfo) error {
	// Check if the employee info already exists
	existingEmployeeInfo, err := r.Get(ctx, tx, data.ID)
	if err != nil {
		return fmt.Errorf("failed to check existing employee info: %w", err)
	}

	if existingEmployeeInfo == nil {
		// If it doesn't exist, create a new record
		return r.Create(ctx, tx, data)
	}

	// If it exists, update the existing record
	if err := tx.Save(data).Error; err != nil {
		return fmt.Errorf("failed to save employee info: %w", err)
	}

	return nil
}
