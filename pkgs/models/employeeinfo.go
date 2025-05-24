package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/brianvoe/gofakeit/v6"
)

////////////////////////////////////////////////////////////////////////////////

type EmployeeInfo struct {
	ID      int64  `gorm:"primaryKey" fake:"-"`
	Name    string `fake:"{firstname}"`
	Age     int    `fake:"{number:20,50}"`
	Address string `fake:"{streetname}"`
	Phone   string `fake:"{phone}"`
	Email   string `fake:"{email}"`

	CreatedAt time.Time      `gorm:"autoCreateTime" fake:"-"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" fake:"-"`
	DeleteAt  gorm.DeletedAt `fake:"-"`
}

func (EmployeeInfo) TableName() string {
	return "employeeinfo"
}

////////////////////////////////////////////////////////////////////////////////

func DummyEmployeeInfo(faker *gofakeit.Faker) *EmployeeInfo {
	var gen EmployeeInfo
	if err := faker.Struct(&gen); err != nil {
		panic(err)
	}

	return &gen
}
