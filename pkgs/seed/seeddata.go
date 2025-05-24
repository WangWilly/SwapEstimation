package seed

import (
	"context"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"github.com/WangWilly/labs-hr-go/pkgs/repos/employeeattendancerepo"
	"github.com/WangWilly/labs-hr-go/pkgs/repos/employeeinforepo"
	"github.com/WangWilly/labs-hr-go/pkgs/repos/employeepositionrepo"
	"github.com/brianvoe/gofakeit/v6"
	"gorm.io/gorm"
)

func SeedData(ctx context.Context, db *gorm.DB) error {
	faker := gofakeit.New(0)

	// Create employee info records
	employeeInfo1 := models.DummyEmployeeInfo(faker)
	employeeInfo2 := models.DummyEmployeeInfo(faker)
	employeeInfo3 := models.DummyEmployeeInfo(faker)
	employeeInfo4 := models.DummyEmployeeInfo(faker)
	employeeInfo5 := models.DummyEmployeeInfo(faker)

	// Create repositories
	employeeInfoRepo := employeeinforepo.New()
	employeePositionRepo := employeepositionrepo.New()
	employeeAttendanceRepo := employeeattendancerepo.New()

	// Insert employee info records
	if err := employeeInfoRepo.Create(ctx, db, employeeInfo1); err != nil {
		return err
	}
	if err := employeeInfoRepo.Create(ctx, db, employeeInfo2); err != nil {
		return err
	}
	if err := employeeInfoRepo.Create(ctx, db, employeeInfo3); err != nil {
		return err
	}
	if err := employeeInfoRepo.Create(ctx, db, employeeInfo4); err != nil {
		return err
	}
	if err := employeeInfoRepo.Create(ctx, db, employeeInfo5); err != nil {
		return err
	}

	// Create employee positions for each employee
	employees := []*models.EmployeeInfo{employeeInfo1, employeeInfo2, employeeInfo3, employeeInfo4, employeeInfo5}
	for _, emp := range employees {
		position := models.DummyEmployeePosition(faker)
		position.EmployeeID = emp.ID

		if err := employeePositionRepo.Create(ctx, db, position, time.Now()); err != nil {
			return err
		}

		// Create attendance records for each employee
		now := time.Now()
		// Create attendance records for the past week
		for i := 0; i < 5; i++ {
			clockInTime := now.AddDate(0, 0, -i).Add(-8 * time.Hour)
			// Add some randomness to clock in time
			clockInTime = clockInTime.Add(time.Duration(faker.IntRange(-30, 30)) * time.Minute)

			_, err := employeeAttendanceRepo.CreateForClockIn(ctx, db, emp.ID, position.ID, clockInTime)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
