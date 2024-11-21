package migrations

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
	"time"
)

const migrationTable = `create table if not exists migrations(
    id text unique primary key,
    created_at integer not null
)`

func Migrate(db *sql.DB) error {
	dir, err := os.ReadDir("./migrations")
	if err != nil {
		return err
	}

	res, err := db.Exec(migrationTable)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected != 0 {
		log.Println("succesfully create migrations table")
	} else {
		log.Println("migration table already exists. skipping creation...")
	}

	tx, err := db.Begin()
	defer tx.Rollback()

	for _, entry := range dir {
		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		rows, err := db.Query("select * from migrations where id = ?", entry.Name())
		defer rows.Close()
		if rows.Next() {
			continue
		}
		file, err := os.ReadFile("./migrations/" + entry.Name())
		if err != nil {
			log.Fatal("could not read migration file: ", err)
		}
		if err != nil {
			log.Fatal("could not create database transaction: ", err)
		}
		_, err = tx.Exec(string(file))
		if err != nil {
			return err
		}
		_, err = tx.Exec("insert into migrations (id, created_at) VALUES (?, ?)", entry.Name(), time.Now().UnixMilli())
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		log.Println("successfully ran migration: ", entry.Name())
	}
	err = tx.Commit()
	log.Println("finished running migrations")
	return nil
}
