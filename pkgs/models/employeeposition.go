package models

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

////////////////////////////////////////////////////////////////////////////////

type EmployeePosition struct {
	ID         int64 `gorm:"primaryKey" fake:"-"`
	EmployeeID int64 `gorm:"index" fake:"{number:1,100}"`

	Position   string  `gorm:"size:100" fake:"{word}"`
	Department string  `gorm:"size:100" fake:"{word}"`
	Salary     float64 `gorm:"type:decimal(10,2)" fake:"{price:1000,5000}"`

	StartDate time.Time `gorm:"type:date" fake:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" fake:"-"`
}

func (EmployeePosition) TableName() string {
	return "employeeposition"
}

////////////////////////////////////////////////////////////////////////////////

func DummyEmployeePosition(faker *gofakeit.Faker) *EmployeePosition {
	var gen EmployeePosition
	if err := faker.Struct(&gen); err != nil {
		panic(err)
	}
	// generate date without time
	gen.StartDate = faker.Date()
	gen.StartDate = time.Date(
		gen.StartDate.Year(),
		gen.StartDate.Month(),
		gen.StartDate.Day(),
		0, 0, 0, 0, time.UTC,
	)

	return &gen
}
