package db

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	db, err := Open()
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()

	err = db.Ping()
	assert.NoError(t, err)

	dbPath := filepath.Join(tmpDir, ".config", "limelight", "limelight.db")
	_, err = os.Stat(dbPath)
	assert.NoError(t, err)
}

func TestOpenCreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	configDir := filepath.Join(tmpDir, ".config", "limelight")
	_, err := os.Stat(configDir)
	assert.True(t, os.IsNotExist(err))

	db, err := Open()
	require.NoError(t, err)
	defer db.Close()

	_, err = os.Stat(configDir)
	assert.NoError(t, err)
}
