package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var schema string = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title TEXT NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);     
CREATE INDEX date_idx ON scheduler (date);`

func Init(dbFile string) error {
	var install bool

	envDBFile := os.Getenv("TODO_DBFILE")
	if len(envDBFile) > 0 {
		dbFile = envDBFile
	}

	_, err := os.Stat(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return err
		}
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

// insert into scheduler (date, title, comment, repeat) values (20260601, kek, lol, rep);
