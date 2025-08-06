package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type status int

const (
	queued status = iota
	processed
	failed
	waiting
	sent
)

var statusString = map[status]string{
	queued:    "queued",
	processed: "processed",
	failed:    "failed",
	waiting:   "waiting",
	sent:      "sent",
}

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
	checkTable := `
  CREATE TABLE IF NOT EXISTS file_process (
    _id INTEGER PRIMARY KEY AUTOINCREMENT,
    filename TEXT NOT NULL,
    status TEXT NOT NULL,
    os_path TEXT NOT NULL
  );`

	_, err = conn.Exec(checkTable)
	if err != nil {
		return err
	}
	return nil
}
