package db

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})

	db, err := Open()
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestRunMigrations(t *testing.T) {
	db := setupTestDB(t)

	err := RunMigrations(db)
	require.NoError(t, err)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	var version string
	err = db.QueryRow("SELECT version FROM schema_migrations").Scan(&version)
	require.NoError(t, err)
	assert.Equal(t, "001_initial_schema.sql", version)
}

func TestRunMigrationsIdempotent(t *testing.T) {
	db := setupTestDB(t)

	err := RunMigrations(db)
	require.NoError(t, err)

	err = RunMigrations(db)
	require.NoError(t, err)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestMigrationsCreateTables(t *testing.T) {
	db := setupTestDB(t)

	err := RunMigrations(db)
	require.NoError(t, err)

	tables := []string{"automations", "triggers", "conditions", "actions", "config"}
	for _, table := range tables {
		var name string
		err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		require.NoError(t, err, "table %s should exist", table)
		assert.Equal(t, table, name)
	}
}

func TestMigrationsForeignKeys(t *testing.T) {
	db := setupTestDB(t)

	err := RunMigrations(db)
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO automations (name, description) VALUES (?, ?)", "test", "test automation")
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO triggers (automation_id, type, config) VALUES (?, ?, ?)", 999, "time", "{}")
	assert.Error(t, err)
}
