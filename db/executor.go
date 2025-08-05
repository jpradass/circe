package db

type status int

const (
	queued status = iota
	processed
	failed
)

var statusString = map[status]string{
	queued:    "queued",
	processed: "processed",
	failed:    "failed",
}

func executorChecks() error {
	checkTable := `
  CREATE TABLE IF NOT EXISTS executor (
    _id INTEGER PRIMARY KEY AUTOINCREMENT,
    filename TEXT NOT NULL,
    status TEXT NOT NULL,
    os_path TEXT NOT NULL,
  );`

	_, err = conn.Exec(checkTable)
	if err != nil {
		return err
	}

	return nil
}
