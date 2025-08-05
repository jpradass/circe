package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var (
	conn *sql.DB
	err  error
)

func Init() error {
	conn, err = sql.Open("sqlite3", "circe.db")
	if err != nil {
		return err
	}

	if err = runChecks(); err != nil {
		return err
	}

	return nil
}

func runChecks() error {
	if err := watcherChecks(); err != nil {
		return err
	}

	if err := executorChecks(); err != nil {
		return err
	}

	return nil
}
