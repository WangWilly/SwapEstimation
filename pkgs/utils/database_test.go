package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDB(t *testing.T) {
	Convey("Given a database configuration", t, func() {
		Convey("When using an unsupported driver", func() {
			cfg := DbConfig{
				Driver: "unsupported",
			}
			db, err := GetDB(cfg)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unsupported database driver: unsupported")
				So(db, ShouldBeNil)
			})
		})

		Convey("When using the MySQL driver with invalid configuration", func() {
			cfg := DbConfig{
				Driver:   "mysql",
				Host:     "nonexistent-host",
				Port:     "3306",
				User:     "invalid",
				Password: "invalid",
				Database: "invalid",
			}
			db, err := GetDB(cfg)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(db, ShouldBeNil)
			})
		})

		Convey("When using the MySQL driver with valid configuration", func() {
			// This test requires an actual MySQL server. Skip if not available
			t.Skip("Skipping test that requires a running MySQL server")

			cfg := DbConfig{
				Driver:   "mysql",
				Host:     "localhost",
				Port:     "3306",
				User:     "root",
				Password: "password",
				Database: "test",
				// Set connection pool settings
				MaxIdleConns:    5,
				MaxOpenConns:    10,
				ConnMaxLifetime: 30,
				ConnMaxIdleTime: 30,
				IsDev:           true,
			}

			db, err := GetDB(cfg)

			Convey("Then it should successfully connect", func() {
				So(err, ShouldBeNil)
				So(db, ShouldNotBeNil)

				// Verify the connection pool settings were applied
				sqlDB, err := db.DB()
				So(err, ShouldBeNil)
				So(sqlDB.Stats().MaxOpenConnections, ShouldEqual, 10)
			})
		})
	})
}

func TestDBContextKey(t *testing.T) {
	Convey("Given the DBContextKey", t, func() {
		Convey("It should have the correct value", func() {
			So(DBContextKey, ShouldEqual, contextKey("gormdb"))
		})
	})
}
