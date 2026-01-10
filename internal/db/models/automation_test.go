package models

import (
	"database/sql"
	"os"
	"testing"

	"github.com/mithilarun/limelight/internal/db"
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

	database, err := db.Open()
	require.NoError(t, err)
	t.Cleanup(func() {
		database.Close()
	})

	err = db.RunMigrations(database)
	require.NoError(t, err)

	return database
}

func TestCreateAutomation(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Morning Lights", "Turn on lights in the morning")
	require.NoError(t, err)
	require.NotNil(t, automation)

	assert.Greater(t, automation.ID, int64(0))
	assert.Equal(t, "Morning Lights", automation.Name)
	assert.Equal(t, "Turn on lights in the morning", automation.Description)
	assert.True(t, automation.Enabled)
	assert.False(t, automation.CreatedAt.IsZero())
	assert.False(t, automation.UpdatedAt.IsZero())
}

func TestCreateAutomationEmptyName(t *testing.T) {
	database := setupTestDB(t)

	_, err := CreateAutomation(database, "", "No name")
	assert.Error(t, err)
}

func TestCreateAutomationDuplicateName(t *testing.T) {
	database := setupTestDB(t)

	_, err := CreateAutomation(database, "Morning Lights", "First")
	require.NoError(t, err)

	_, err = CreateAutomation(database, "Morning Lights", "Second")
	assert.Error(t, err)
}

func TestGetAutomation(t *testing.T) {
	database := setupTestDB(t)

	created, err := CreateAutomation(database, "Test Automation", "Test description")
	require.NoError(t, err)

	retrieved, err := GetAutomation(database, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Name, retrieved.Name)
	assert.Equal(t, created.Description, retrieved.Description)
}

func TestGetAutomationNotFound(t *testing.T) {
	database := setupTestDB(t)

	_, err := GetAutomation(database, 999)
	assert.Error(t, err)
}

func TestListAutomations(t *testing.T) {
	database := setupTestDB(t)

	_, err := CreateAutomation(database, "First", "First automation")
	require.NoError(t, err)

	_, err = CreateAutomation(database, "Second", "Second automation")
	require.NoError(t, err)

	automations, err := ListAutomations(database)
	require.NoError(t, err)
	assert.Len(t, automations, 2)
}

func TestListEnabledAutomations(t *testing.T) {
	database := setupTestDB(t)

	first, err := CreateAutomation(database, "First", "First automation")
	require.NoError(t, err)

	_, err = CreateAutomation(database, "Second", "Second automation")
	require.NoError(t, err)

	err = SetEnabled(database, first.ID, false)
	require.NoError(t, err)

	automations, err := ListEnabledAutomations(database)
	require.NoError(t, err)
	assert.Len(t, automations, 1)
	assert.Equal(t, "Second", automations[0].Name)
}

func TestUpdateAutomation(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Original", "Original description")
	require.NoError(t, err)

	err = UpdateAutomation(database, automation.ID, "Updated", "Updated description")
	require.NoError(t, err)

	updated, err := GetAutomation(database, automation.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated", updated.Name)
	assert.Equal(t, "Updated description", updated.Description)
}

func TestUpdateAutomationEmptyName(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Original", "Original description")
	require.NoError(t, err)

	err = UpdateAutomation(database, automation.ID, "", "Updated description")
	assert.Error(t, err)
}

func TestUpdateAutomationNotFound(t *testing.T) {
	database := setupTestDB(t)

	err := UpdateAutomation(database, 999, "Updated", "Updated description")
	assert.Error(t, err)
}

func TestSetEnabled(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)
	assert.True(t, automation.Enabled)

	err = SetEnabled(database, automation.ID, false)
	require.NoError(t, err)

	updated, err := GetAutomation(database, automation.ID)
	require.NoError(t, err)
	assert.False(t, updated.Enabled)

	err = SetEnabled(database, automation.ID, true)
	require.NoError(t, err)

	updated, err = GetAutomation(database, automation.ID)
	require.NoError(t, err)
	assert.True(t, updated.Enabled)
}

func TestSetEnabledNotFound(t *testing.T) {
	database := setupTestDB(t)

	err := SetEnabled(database, 999, false)
	assert.Error(t, err)
}

func TestDeleteAutomation(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	err = DeleteAutomation(database, automation.ID)
	require.NoError(t, err)

	_, err = GetAutomation(database, automation.ID)
	assert.Error(t, err)
}

func TestDeleteAutomationNotFound(t *testing.T) {
	database := setupTestDB(t)

	err := DeleteAutomation(database, 999)
	assert.Error(t, err)
}

func TestDeleteAutomationCascade(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	_, err = CreateTrigger(database, automation.ID, TriggerTypeTime, map[string]interface{}{"hour": 9, "minute": 0})
	require.NoError(t, err)

	_, err = CreateCondition(database, automation.ID, ConditionTypeWeekday, map[string]interface{}{})
	require.NoError(t, err)

	_, err = CreateAction(database, automation.ID, ActionTypeScene, map[string]interface{}{"scene_id": "abc"}, 0)
	require.NoError(t, err)

	err = DeleteAutomation(database, automation.ID)
	require.NoError(t, err)

	triggers, err := GetTriggers(database, automation.ID)
	require.NoError(t, err)
	assert.Len(t, triggers, 0)

	conditions, err := GetConditions(database, automation.ID)
	require.NoError(t, err)
	assert.Len(t, conditions, 0)

	actions, err := GetActions(database, automation.ID)
	require.NoError(t, err)
	assert.Len(t, actions, 0)
}
