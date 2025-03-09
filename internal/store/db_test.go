package store

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"os"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	dbTest, err := sql.Open("sqlite3", ":memory:")
	db = dbTest
	if err != nil {
		return -1, fmt.Errorf("could not connect to database: %w", err)
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	migration, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations",
		"sqlite3", driver)
	if err != nil {
		return -2, fmt.Errorf("could not create migration object: %w", err)
	}
	err = migration.Up()
	if err != nil {
		return -3, fmt.Errorf("could not migrate database: %w", err)
	}

	defer func() {
		for _, t := range []string{"users", "groups", "groups_users"} {
			_, _ = db.Exec("DELETE FROM ?", t)
		}
		db.Close()
	}()
	return m.Run(), nil
}
