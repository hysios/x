package repos

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func openDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	columns := []string{"version"}
	mock.ExpectQuery("SELECT VERSION()").WithArgs().WillReturnRows(
		mock.NewRows(columns).FromCSVString("1"),
	)

	gdb, err := gorm.Open(mysql.New(
		mysql.Config{
			Conn:       db,
			DriverName: "mysql",
		}), &gorm.Config{})
	_ = err

	return gdb, mock
}
