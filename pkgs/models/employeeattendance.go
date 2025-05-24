package models

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

////////////////////////////////////////////////////////////////////////////////

type EmployeeAttendance struct {
	ID         int64 `gorm:"primaryKey" fake:"-"`
	EmployeeID int64 `gorm:"index" fake:"{number:1,100}"`

	PositionID int64     `gorm:"index" fake:"{number:1,100}"`
	ClockIn    time.Time `gorm:"datetime" fake:"-"`
	ClockOut   time.Time `gorm:"datetime" fake:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" fake:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" fake:"-"`
}

func (EmployeeAttendance) TableName() string {
	return "employeeattendance"
}

////////////////////////////////////////////////////////////////////////////////

func DummyEmployeeAttendance(faker *gofakeit.Faker) *EmployeeAttendance {
	var gen EmployeeAttendance
	if err := faker.Struct(&gen); err != nil {
		panic(err)
	}

	gen.ClockIn = faker.Date()
	gen.ClockOut = gen.ClockIn.Add(time.Hour * 8)

	return &gen
}
