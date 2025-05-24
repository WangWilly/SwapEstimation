package attendance

import (
	"net/http"
	"strconv"

	"github.com/WangWilly/labs-hr-go/pkgs/dtos"
	"github.com/WangWilly/labs-hr-go/pkgs/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

////////////////////////////////////////////////////////////////////////////////

func (c *Controller) Get(ctx *gin.Context) {
	logger := log.Ctx(ctx.Request.Context())

	employeeID, ok := ctx.Params.Get("employee_id")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	employeeIDInt, err := strconv.ParseInt(employeeID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	////////////////////////////////////////////////////////////////////////////

	cached, err := c.cacheManager.GetAttendanceV1(ctx, employeeIDInt)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get attendance from cache")
	}
	if cached != nil {
		logger.Info().Msg("Cache hit")
		ctx.JSON(http.StatusOK, cached)
		return
	}

	////////////////////////////////////////////////////////////////////////////

	// Get the last attendance record of the employee
	currAttendance, err := c.employeeAttendanceRepo.Last(ctx, c.db, employeeIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get last attendance"})
		return
	}
	if currAttendance == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "attendance not found"})
		return
	}

	clockOutTime := ""
	if currAttendance.ClockIn != currAttendance.ClockOut {
		// If the clock-in and clock-out times are the same, it means the employee has not clocked out yet
		clockOutTime = utils.FormatedTime(currAttendance.ClockOut)
	}

	////////////////////////////////////////////////////////////////////////////
	resp := dtos.AttendanceV1Response{
		AttendanceID: currAttendance.ID,
		PositionID:   currAttendance.PositionID,
		ClockInTime:  utils.FormatedTime(currAttendance.ClockIn),
		ClockOutTime: clockOutTime,
	}

	// Cache the attendance record
	if err := c.cacheManager.SetAttendanceV1(ctx, employeeIDInt, resp, 0); err != nil {
		logger.Error().Err(err).Msg("Failed to set attendance to cache")
	}

	// Return the attendance record
	ctx.JSON(http.StatusOK, resp)
}
