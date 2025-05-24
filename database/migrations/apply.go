package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

////////////////////////////////////////////////////////////////////////////////

var migrationList = []*gormigrate.Migration{
	m00001,
}

func Apply(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrationList)
	if err := m.Migrate(); err != nil {
		return err
	}

	return nil
}
