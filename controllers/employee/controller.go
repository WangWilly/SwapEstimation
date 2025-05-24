package employee

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

////////////////////////////////////////////////////////////////////////////////

type Config struct {
}

type Controller struct {
	cfg Config
	db  *gorm.DB

	timeModule           TimeModule
	employeeInfoRepo     EmployeeInfoRepo
	employeePositionRepo EmployeePositionRepo
	cacheManager         CacheManager
}

func NewController(
	cfg Config,
	db *gorm.DB,
	timeModule TimeModule,
	employeeInfoRepo EmployeeInfoRepo,
	employeePositionRepo EmployeePositionRepo,
	cacheManager CacheManager,
) *Controller {
	return &Controller{
		cfg:                  cfg,
		db:                   db,
		timeModule:           timeModule,
		employeeInfoRepo:     employeeInfoRepo,
		employeePositionRepo: employeePositionRepo,
		cacheManager:         cacheManager,
	}
}

func (c *Controller) RegisterRoutes(r *gin.Engine) {
	////////////////////////////////////////////////////////////////////////////
	// employee management
	r.POST("/employee", c.Create)
	r.GET("/employee/:id", c.Get)
	// r.GET("/employee", c.List)
	r.PUT("/employee/:id", c.Update)
	// r.DELETE("/employee/:id", c.Delete)

	////////////////////////////////////////////////////////////////////////////
	// position management
	r.POST("/promote/:id", c.Promote)
}
