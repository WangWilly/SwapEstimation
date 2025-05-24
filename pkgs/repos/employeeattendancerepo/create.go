package employeeattendancerepo

import (
	"context"
	"fmt"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"gorm.io/gorm"
)

func (r *repo) CreateForClockIn(ctx context.Context, tx *gorm.DB, employeeID int64, positionID int64, clockInTime time.Time) (*models.EmployeeAttendance, error) {
	// Create a new employee attendance record
	employeeAttendance := &models.EmployeeAttendance{
		EmployeeID: employeeID,
		PositionID: positionID,
		ClockIn:    clockInTime,
		ClockOut:   clockInTime, // Initialize ClockOut to the same time as ClockIn
	}

	// Execute the insert query
	if err := tx.Create(employeeAttendance).Error; err != nil {
		return nil, fmt.Errorf("failed to create employee attendance: %w", err)
	}

	// Return the created record
	return employeeAttendance, nil
}
