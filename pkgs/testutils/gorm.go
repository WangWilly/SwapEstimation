package testutils

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/WangWilly/swap-estimation/pkgs/utils"
	"github.com/sethvargo/go-envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

////////////////////////////////////////////////////////////////////////////////

func tableName(tx *gorm.DB, model interface{}) string {
	stmt := &gorm.Statement{DB: tx}
	stmt.Parse(model)
	return stmt.Schema.Table
}

func MustClearTable(t *testing.T, tx *gorm.DB, model interface{}) {
	tableName := tableName(tx, model)
	_, err := tx.Raw(fmt.Sprintf("DELETE FROM %s", tableName)).Rows()
	if err != nil {
		t.Fatal(err)
	}
}

////////////////////////////////////////////////////////////////////////////////

func GetMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	ctx := t.Context()

	////////////////////////////////////////////////////////////////////////////

	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	// Add expectation for the GORM initialization query
	mockDB.ExpectQuery("SELECT VERSION()").WillReturnRows(
		sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7.32"))

	////////////////////////////////////////////////////////////////////////////

	dbCfg := utils.DbConfig{}
	if err := envconfig.Process(ctx, &dbCfg); err != nil {
		t.Fatal(err)
	}

	var dialector gorm.Dialector
	switch dbCfg.Driver {
	case "mysql":
		dialector = mysql.New(mysql.Config{Conn: db})
	default:
		t.Fatalf("unsupported driver: %s", dbCfg.Driver)
	}

	////////////////////////////////////////////////////////////////////////////

	gormDB, err := gorm.Open(
		dialector,
		&gorm.Config{SkipDefaultTransaction: true},
	)
	if err != nil {
		t.Fatal(err)
	}

	////////////////////////////////////////////////////////////////////////////

	gormDB.Logger.LogMode(0)
	return gormDB, mockDB
}
