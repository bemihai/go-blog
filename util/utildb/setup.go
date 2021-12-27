package utildb

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExecDir executes all sql files found in the given folder.
func ExecDir(db *sql.DB, dir string) error {
	paths, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, f := range paths {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			err := ExecFile(db, filepath.Join(dir, f.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ExecFile executes the given sql file into the given db.
func ExecFile(db *sql.DB, file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(data))
	return err
}

// ExecFiles executes the given sql files into the given db.
func ExecFiles(db *sql.DB, files []string) error {
	for _, file := range files {
		err := ExecFile(db, file)
		if err != nil {
			return fmt.Errorf("exec %s: %w", file, err)
		}
	}
	return nil
}
