package employee

import (
	"strconv"

	"github.com/WangWilly/labs-hr-go/pkgs/dtos"
	"github.com/WangWilly/labs-hr-go/pkgs/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

////////////////////////////////////////////////////////////////////////////////

func (c *Controller) Get(ctx *gin.Context) {
	logger := log.Ctx(ctx.Request.Context())

	id, ok := ctx.Params.Get("id")
	if !ok {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	// Convert id to int64
	employeeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	////////////////////////////////////////////////////////////////////////////

	// Check if the employee detail is in cache
	cacheData, err := c.cacheManager.GetEmployeeDetailV1(ctx, employeeID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get employee detail from cache")
	}
	if err == nil && cacheData != nil {
		logger.Info().Msg("Cache hit")
		ctx.JSON(200, cacheData)
		return
	}

	////////////////////////////////////////////////////////////////////////////

	employeeInfo, err := c.employeeInfoRepo.MustGet(ctx, c.db, employeeID)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "employee not found"})
		return
	}

	nowTime := c.timeModule.Now()
	employeePosition, err := c.employeePositionRepo.GetCurrentByEmployeeID(ctx, c.db, employeeID, nowTime)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "employee position not found"})
		return
	}

	////////////////////////////////////////////////////////////////////////////

	response := &dtos.EmployeeV1Response{
		EmployeeID: employeeInfo.ID,
		Name:       employeeInfo.Name,
		Age:        employeeInfo.Age,
		Phone:      employeeInfo.Phone,
		Email:      employeeInfo.Email,
		Address:    employeeInfo.Address,
		CreatedAt:  utils.FormatedTime(employeeInfo.CreatedAt),
		UpdatedAt:  utils.FormatedTime(employeeInfo.UpdatedAt),

		PositionID: employeePosition.ID,
		Position:   employeePosition.Position,
		Department: employeePosition.Department,
		Salary:     employeePosition.Salary,
		StartDate:  utils.FormatedTime(employeePosition.StartDate),
	}

	// Cache the employee detail
	if err := c.cacheManager.SetEmployeeDetailV1(ctx, employeeID, *response, 0); err != nil {
		logger.Error().Err(err).Msg("Failed to cache employee detail")
	}

	ctx.JSON(200, response)
}
