package attendance

import (
	"fmt"
	"net/http"

	"github.com/WangWilly/labs-hr-go/pkgs/dtos"
	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"github.com/WangWilly/labs-hr-go/pkgs/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

////////////////////////////////////////////////////////////////////////////////

type CreateRequest struct {
	EmployeeID int64 `json:"employee_id" binding:"required"`
}

////////////////////////////////////////////////////////////////////////////////

func (c *Controller) Create(ctx *gin.Context) {
	logger := log.Ctx(ctx.Request.Context())

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	////////////////////////////////////////////////////////////////////////////

	// Get the employee's current position
	positionID, err := c.getEmployeePosition(ctx, req.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get employee position"})
		return
	}

	// Get the employee's current attendance
	attendance, err := c.getEmployeeAttendance(ctx, req.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get employee attendance"})
		return
	}

	////////////////////////////////////////////////////////////////////////////

	// Create or update the attendance record
	attendanceResponse, err := c.createClockIn(ctx, req.EmployeeID, positionID, attendance)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create/update attendance"})
		return
	}
	// Cache the attendance record
	if err := c.cacheManager.SetAttendanceV1(ctx, req.EmployeeID, *attendanceResponse, 0); err != nil {
		logger.Error().Err(err).Msg("Failed to cache attendance")
	}

	ctx.JSON(http.StatusCreated, attendanceResponse)
}

////////////////////////////////////////////////////////////////////////////////

func (c *Controller) getEmployeePosition(ctx *gin.Context, employeeID int64) (int64, error) {
	logger := log.Ctx(ctx.Request.Context())

	if employeeID <= 0 {
		return 0, fmt.Errorf("invalid employee ID")
	}

	employeePosition, err := c.cacheManager.GetEmployeeDetailV1(ctx, employeeID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get employee position from cache")
	}
	if employeePosition != nil {
		return employeePosition.PositionID, nil
	}

	// Get the current position of the employee
	dbEmployeePosition, err := c.employeePositionRepo.GetCurrentByEmployeeID(ctx, c.db, employeeID, c.timeModule.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to get employee position: %w", err)
	}
	if dbEmployeePosition == nil {
		return 0, fmt.Errorf("employee position not found")
	}
	return dbEmployeePosition.ID, nil
}

func (c *Controller) getEmployeeAttendance(ctx *gin.Context, employeeID int64) (*models.EmployeeAttendance, error) {
	if employeeID <= 0 {
		return nil, fmt.Errorf("invalid employee ID")
	}

	// Get the current attendance of the employee
	dbAttendance, err := c.employeeAttendanceRepo.Last(ctx, c.db, employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee attendance: %w", err)
	}
	if dbAttendance == nil {
		return nil, nil
	}
	return dbAttendance, nil
}

func (c *Controller) createClockIn(
	ctx *gin.Context,
	employeeID int64,
	positionID int64,
	currAttendance *models.EmployeeAttendance,
) (*dtos.AttendanceV1Response, error) {
	if employeeID <= 0 {
		return nil, fmt.Errorf("invalid employee ID")
	}
	if positionID <= 0 {
		return nil, fmt.Errorf("invalid position ID")
	}

	if currAttendance == nil || currAttendance.ClockIn != currAttendance.ClockOut {
		// Create a new attendance record for clock-in
		newAttendance, err := c.employeeAttendanceRepo.CreateForClockIn(
			ctx,
			c.db,
			employeeID,
			positionID,
			c.timeModule.Now(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create attendance: %w", err)
		}
		return &dtos.AttendanceV1Response{
			AttendanceID: newAttendance.ID,
			PositionID:   newAttendance.PositionID,
			ClockInTime:  utils.FormatedTime(newAttendance.ClockIn),
			ClockOutTime: "",
		}, nil
	}

	////////////////////////////////////////////////////////////////////////////

	currAttendance, err := c.employeeAttendanceRepo.UpdateForClockOut(
		ctx,
		c.db,
		currAttendance.ID,
		c.timeModule.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update attendance: %w", err)
	}
	return &dtos.AttendanceV1Response{
		AttendanceID: currAttendance.ID,
		PositionID:   currAttendance.PositionID,
		ClockInTime:  utils.FormatedTime(currAttendance.ClockIn),
		ClockOutTime: utils.FormatedTime(currAttendance.ClockOut),
	}, nil
}
