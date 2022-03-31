package migrate

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"io/fs"
)

func InitAutoMigrateForMysqlFunc(dsn string, f fs.FS) func() error {
	return func() (err error) {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return
		}
		defer db.Close()

		instance, err := mysql.WithInstance(db, &mysql.Config{})
		if err != nil {
			return
		}

		driver := NewFsDriver()
		driver.Add(f)

		m, err := migrate.NewWithInstance("", driver, "main", instance)
		if err != nil {
			return
		}

		err = m.Up()
		if err != nil && err == migrate.ErrNoChange {
			err = nil
		}
		return
	}
}
