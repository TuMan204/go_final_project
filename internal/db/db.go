package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const schema string = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title TEXT NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);     
CREATE INDEX date_idx ON scheduler (date);`

func Init() (*sql.DB, error) {
	var install bool

	dbFile := "scheduler.db"
	envDBFile := os.Getenv("TODO_DBFILE")
	if len(envDBFile) > 0 {
		dbFile = envDBFile
	}
	dbPath := filepath.Join("./", dbFile)

	_, err := os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
