package db

import "fmt"

func AddFile(filename, os_path string) error {
	fmt.Printf("adding file %s with path %s to db\n", filename, os_path)

	insertSQL := `INSERT INTO file_process (filename, status, os_path) VALUES (?,?,?)`
	statement, err := conn.Prepare(insertSQL)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(filename, waiting, os_path)
	if err != nil {
		return err
	}

	return nil
}

func SetStatusToFile(filename string) error {
	fmt.Printf("setting `sent` to %s\n", filename)

	setStatusStmt := `UPDATE file_process SET status = ? WHERE filename = ?`

	_, err := conn.Exec(setStatusStmt, sent, filename)
	if err != nil {
		return err
	}

	return nil
}
