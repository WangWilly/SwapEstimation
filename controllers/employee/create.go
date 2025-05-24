package employee

import (
	"net/http"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/dtos"
	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"github.com/WangWilly/labs-hr-go/pkgs/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

////////////////////////////////////////////////////////////////////////////////

type CreateRequest struct {
	Name    string `json:"name"    binding:"required"`
	Age     int    `json:"age"     binding:"required"`
	Address string `json:"address" binding:"required"`
	Phone   string `json:"phone"   binding:"required"`
	Email   string `json:"email"   binding:"required"`

	Position   string  `json:"position"   binding:"required"`
	Department string  `json:"department" binding:"required"`
	Salary     float64 `json:"salary"     binding:"required"`
	StartDate  int64   `json:"start_date" binding:"required"`
}

type CreateResponse struct {
	EmployeeID int64 `json:"employee_id"`
	PositionID int64 `json:"position_id"`
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

	// Create the employee info
	employeeInfo := &models.EmployeeInfo{
		Name:    req.Name,
		Age:     req.Age,
		Address: req.Address,
		Phone:   req.Phone,
		Email:   req.Email,
	}
	if err := c.employeeInfoRepo.Create(ctx, c.db, employeeInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create employee info"})
		return
	}

	// Create the employee position
	employeePosition := &models.EmployeePosition{
		EmployeeID: employeeInfo.ID,
		Position:   req.Position,
		Department: req.Department,
		Salary:     req.Salary,
		StartDate:  time.Unix(req.StartDate, 0),
	}
	nowTime := c.timeModule.Now()
	if err := c.employeePositionRepo.Create(ctx, c.db, employeePosition, nowTime); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create employee position"})
		return
	}

	////////////////////////////////////////////////////////////////////////////

	// Cache the employee detail
	employeeDetail := dtos.EmployeeV1Response{
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
	if err := c.cacheManager.SetEmployeeDetailV1(ctx, employeeInfo.ID, employeeDetail, 0); err != nil {
		logger.Error().Err(err).Msg("Failed to cache employee detail")
	}

	ctx.JSON(http.StatusCreated, CreateResponse{
		EmployeeID: employeeInfo.ID,
		PositionID: employeePosition.ID,
	})
}
