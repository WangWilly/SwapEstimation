package attendance

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

	timeModule             TimeModule
	employeePositionRepo   EmployeePositionRepo
	employeeAttendanceRepo EmployeeAttendanceRepo
	cacheManager           CacheManager
}

func NewController(
	cfg Config,
	db *gorm.DB,
	timeModule TimeModule,
	employeePositionRepo EmployeePositionRepo,
	employeeAttendanceRepo EmployeeAttendanceRepo,
	cacheManage CacheManager,
) *Controller {
	return &Controller{
		cfg:                    cfg,
		db:                     db,
		timeModule:             timeModule,
		employeePositionRepo:   employeePositionRepo,
		employeeAttendanceRepo: employeeAttendanceRepo,
		cacheManager:           cacheManage,
	}
}

func (c *Controller) RegisterRoutes(r *gin.Engine) {
	////////////////////////////////////////////////////////////////////////////
	// attendance management
	r.POST("/attendance", c.Create)
	r.GET("/attendance/:employee_id", c.Get)
}
