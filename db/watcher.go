package db

func watcherChecks() error {
	checkTable := `
  CREATE TABLE IF NOT EXISTS watcher (
    _id INTEGER PRIMARY KEY AUTOINCREMENT,
    filename TEXT NOT NULL,
    os_path TEXT NOT NULL,
  );`

	_, err = conn.Exec(checkTable)
	if err != nil {
		return err
	}

	return nil
}
