package employeeattendancerepo

import (
	"context"
	"fmt"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"gorm.io/gorm"
)

func (r *repo) UpdateForClockOut(ctx context.Context, tx *gorm.DB, attendanceID int64, clockOutTime time.Time) (*models.EmployeeAttendance, error) {
	// Create a variable to hold the result
	var employeeAttendance models.EmployeeAttendance

	// Execute the query to find the attendance record by ID
	if err := tx.Where("id = ?", attendanceID).First(&employeeAttendance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get employee attendance: %w", err)
	}
	// Check if the attendance record is already clocked out
	if employeeAttendance.ClockOut != employeeAttendance.ClockIn {
		return nil, fmt.Errorf("attendance record already clocked out")
	}

	// Update the clock-out time
	employeeAttendance.ClockOut = clockOutTime
	// Save the updated record
	if err := tx.Save(&employeeAttendance).Error; err != nil {
		return nil, fmt.Errorf("failed to update employee attendance: %w", err)
	}
	// Return the updated record
	return &employeeAttendance, nil
}
