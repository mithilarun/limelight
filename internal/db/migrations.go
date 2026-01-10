package db

import (
	"database/sql"
	"embed"
	"fmt"
	"sort"

	"github.com/cockroachdb/errors"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// RunMigrations applies all pending database migrations
func RunMigrations(db *sql.DB) error {
	if err := createMigrationsTable(db); err != nil {
		return errors.Wrap(err, "failed to create migrations table")
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return errors.Wrap(err, "failed to get applied migrations")
	}

	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return errors.Wrap(err, "failed to read migrations directory")
	}

	var migrationNames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			migrationNames = append(migrationNames, entry.Name())
		}
	}
	sort.Strings(migrationNames)

	for _, name := range migrationNames {
		if applied[name] {
			continue
		}

		if err := applyMigration(db, name); err != nil {
			return errors.Wrapf(err, "failed to apply migration %s", name)
		}
	}

	return nil
}

// createMigrationsTable creates the schema_migrations table if it doesn't exist
func createMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return errors.Wrap(err, "failed to create schema_migrations table")
}

// getAppliedMigrations returns a set of already applied migration versions
func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, errors.Wrap(err, "failed to query schema_migrations")
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, errors.Wrap(err, "failed to scan migration version")
		}
		applied[version] = true
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating migration rows")
	}

	return applied, nil
}

// applyMigration applies a single migration within a transaction
func applyMigration(db *sql.DB, name string) error {
	content, err := migrationFiles.ReadFile(fmt.Sprintf("migrations/%s", name))
	if err != nil {
		return errors.Wrap(err, "failed to read migration file")
	}

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(content)); err != nil {
		return errors.Wrap(err, "failed to execute migration SQL")
	}

	if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", name); err != nil {
		return errors.Wrap(err, "failed to record migration")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit migration transaction")
	}

	return nil
}
