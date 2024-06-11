package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./5ch.db")
	if err != nil {
		log.Fatal(err)
	}

	createThreadTableQuery := `
    CREATE TABLE IF NOT EXISTS threads (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err = DB.Exec(createThreadTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	createReplyTableQuery := `
    CREATE TABLE IF NOT EXISTS replies (
        id INTEGER NOT NULL,
        thread_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (id, thread_id),
        FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE
    );
    `
	_, err = DB.Exec(createReplyTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}
