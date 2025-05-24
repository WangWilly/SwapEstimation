package employee

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

////////////////////////////////////////////////////////////////////////////////

type UpdateRequest struct {
	Name    string `json:"name"   `
	Age     int    `json:"age"    `
	Address string `json:"address"`
	Phone   string `json:"phone"  `
	Email   string `json:"email"  `
}

type UpdateResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

func (c *Controller) Update(ctx *gin.Context) {
	logger := log.Ctx(ctx.Request.Context())

	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")
	// Convert id to int64
	employeeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	////////////////////////////////////////////////////////////////////////////

	employeeInfo, err := c.employeeInfoRepo.MustGet(ctx, c.db, employeeID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
		return
	}

	if req.Name != "" {
		employeeInfo.Name = req.Name
	}
	if req.Age != 0 {
		employeeInfo.Age = req.Age
	}
	if req.Address != "" {
		employeeInfo.Address = req.Address
	}
	if req.Phone != "" {
		employeeInfo.Phone = req.Phone
	}
	if req.Email != "" {
		employeeInfo.Email = req.Email
	}

	if err := c.employeeInfoRepo.Save(ctx, c.db, employeeInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update employee info"})
		return
	}

	////////////////////////////////////////////////////////////////////////////

	// Update the cache
	employeeDetail, err := c.cacheManager.GetEmployeeDetailV1(ctx, employeeID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get employee detail from cache")
	}
	if err == nil && employeeDetail != nil {
		employeeDetail.Name = employeeInfo.Name
		employeeDetail.Age = employeeInfo.Age
		employeeDetail.Address = employeeInfo.Address
		employeeDetail.Phone = employeeInfo.Phone
		employeeDetail.Email = employeeInfo.Email

		if err := c.cacheManager.SetEmployeeDetailV1(ctx, employeeID, *employeeDetail, 0); err != nil {
			logger.Error().Err(err).Msg("Failed to cache employee detail")
		}
	}

	ctx.JSON(http.StatusOK, UpdateResponse{
		ID:      employeeInfo.ID,
		Name:    employeeInfo.Name,
		Age:     employeeInfo.Age,
		Address: employeeInfo.Address,
		Phone:   employeeInfo.Phone,
		Email:   employeeInfo.Email,
	})
}
