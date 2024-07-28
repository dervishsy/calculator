package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	*sql.DB
}

func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create tables if not exist
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS expressions (
            id TEXT PRIMARY KEY,
            expression TEXT,
            status TEXT,
            result REAL
        );
        CREATE TABLE IF NOT EXISTS tasks (
            id TEXT PRIMARY KEY,
            expr_id TEXT,
            arg_left TEXT,
            arg_right TEXT,
            operation TEXT,
            result REAL
        );
        CREATE TABLE IF NOT EXISTS sent_tasks (
            task_id TEXT PRIMARY KEY
        );
        CREATE TABLE IF NOT EXISTS task_owners (
            child_id TEXT,
            parent_id TEXT,
            PRIMARY KEY (child_id, parent_id)
        );
        CREATE TABLE IF NOT EXISTS expressions_root (
            task_id TEXT PRIMARY KEY,
            expr_id TEXT
        );
    `)
	if err != nil {
		return nil, err
	}

	return &SQLiteDB{DB: db}, nil
}
