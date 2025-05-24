package employeeinforepo

import (
	"context"
	"fmt"

	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"gorm.io/gorm"
)

func (r *repo) Create(ctx context.Context, tx *gorm.DB, data *models.EmployeeInfo) error {
	if err := tx.
		Create(data).Error; err != nil {
		return fmt.Errorf("failed to create employee info: %w", err)
	}

	return nil
}
