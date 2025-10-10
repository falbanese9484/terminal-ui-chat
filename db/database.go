package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

type Database struct {
	db   *sql.DB
	path string
}

func NewDataBase(dbPath string) (*Database, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	dsn := fmt.Sprintf(DATABASE_PATH, dbPath)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	database := &Database{
		db:   db,
		path: dbPath,
	}

	//if err := database.migrate(); err != nil {
	//	return nil, fmt.Errorf("migration failed: %w", err)
	//}

	return database, nil
}
