package db

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	_ "github.com/mattn/go-sqlite3"
)

// Open opens a connection to the SQLite database.
// The database file is created in ~/.config/limelight/limelight.db
func Open() (*sql.DB, error) {
	dbPath, err := getDatabasePath()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get database path")
	}

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}

	// Configure connection pool
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, errors.Wrap(err, "failed to ping database")
	}

	return db, nil
}

// getDatabasePath returns the path to the database file
func getDatabasePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get home directory")
	}

	configDir := filepath.Join(homeDir, ".config", "limelight")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", errors.Wrap(err, "failed to create config directory")
	}

	return filepath.Join(configDir, "limelight.db"), nil
}
