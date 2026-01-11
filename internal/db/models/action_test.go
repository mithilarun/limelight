package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAction(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	config := map[string]interface{}{
		"scene_id": "abc123",
	}

	action, err := CreateAction(database, automation.ID, ActionTypeScene, config, 0)
	require.NoError(t, err)
	require.NotNil(t, action)

	assert.Greater(t, action.ID, int64(0))
	assert.Equal(t, automation.ID, action.AutomationID)
	assert.Equal(t, ActionTypeScene, action.Type)
	assert.Equal(t, 0, action.OrderIndex)

	var parsed map[string]interface{}
	err = json.Unmarshal(action.Config, &parsed)
	require.NoError(t, err)
	assert.Equal(t, "abc123", parsed["scene_id"])
}

func TestCreateActionInvalidType(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	_, err = CreateAction(database, automation.ID, ActionType("invalid"), map[string]interface{}{}, 0)
	assert.Error(t, err)
}

func TestGetActions(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	_, err = CreateAction(database, automation.ID, ActionTypeScene, map[string]interface{}{"scene_id": "first"}, 1)
	require.NoError(t, err)

	_, err = CreateAction(database, automation.ID, ActionTypeLight, map[string]interface{}{"light_id": "second"}, 0)
	require.NoError(t, err)

	actions, err := GetActions(database, automation.ID)
	require.NoError(t, err)
	assert.Len(t, actions, 2)

	assert.Equal(t, 0, actions[0].OrderIndex)
	assert.Equal(t, ActionTypeLight, actions[0].Type)

	assert.Equal(t, 1, actions[1].OrderIndex)
	assert.Equal(t, ActionTypeScene, actions[1].Type)
}

func TestGetActionsEmpty(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	actions, err := GetActions(database, automation.ID)
	require.NoError(t, err)
	assert.Len(t, actions, 0)
}

func TestDeleteAction(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	action, err := CreateAction(database, automation.ID, ActionTypeScene, map[string]interface{}{"scene_id": "abc"}, 0)
	require.NoError(t, err)

	err = DeleteAction(database, action.ID)
	require.NoError(t, err)

	actions, err := GetActions(database, automation.ID)
	require.NoError(t, err)
	assert.Len(t, actions, 0)
}

func TestDeleteActionNotFound(t *testing.T) {
	database := setupTestDB(t)

	err := DeleteAction(database, 999)
	assert.Error(t, err)
}

func TestCreateActionNegativeOrderIndex(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	_, err = CreateAction(database, automation.ID, ActionTypeScene, map[string]interface{}{"scene_id": "abc"}, -1)
	assert.Error(t, err)
}

func TestValidateActionType(t *testing.T) {
	testCases := []struct {
		name        string
		actionType  ActionType
		expectError bool
	}{
		{
			name:        "valid light action",
			actionType:  ActionTypeLight,
			expectError: false,
		},
		{
			name:        "valid scene action",
			actionType:  ActionTypeScene,
			expectError: false,
		},
		{
			name:        "valid group action",
			actionType:  ActionTypeGroup,
			expectError: false,
		},
		{
			name:        "invalid action type",
			actionType:  ActionType("invalid"),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateActionType(tc.actionType)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
