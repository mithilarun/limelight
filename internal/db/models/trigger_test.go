package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTrigger(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	config := map[string]interface{}{
		"hour":   9,
		"minute": 0,
	}

	trigger, err := CreateTrigger(database, automation.ID, TriggerTypeTime, config)
	require.NoError(t, err)
	require.NotNil(t, trigger)

	assert.Greater(t, trigger.ID, int64(0))
	assert.Equal(t, automation.ID, trigger.AutomationID)
	assert.Equal(t, TriggerTypeTime, trigger.Type)

	var parsed map[string]interface{}
	err = json.Unmarshal(trigger.Config, &parsed)
	require.NoError(t, err)
	assert.Equal(t, float64(9), parsed["hour"])
	assert.Equal(t, float64(0), parsed["minute"])
}

func TestCreateTriggerInvalidType(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	_, err = CreateTrigger(database, automation.ID, TriggerType("invalid"), map[string]interface{}{})
	assert.Error(t, err)
}

func TestGetTriggers(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	_, err = CreateTrigger(database, automation.ID, TriggerTypeTime, map[string]interface{}{"hour": 9})
	require.NoError(t, err)

	_, err = CreateTrigger(database, automation.ID, TriggerTypeSunrise, map[string]interface{}{"offset_minutes": 0})
	require.NoError(t, err)

	triggers, err := GetTriggers(database, automation.ID)
	require.NoError(t, err)
	assert.Len(t, triggers, 2)
}

func TestGetTriggersEmpty(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	triggers, err := GetTriggers(database, automation.ID)
	require.NoError(t, err)
	assert.Len(t, triggers, 0)
}

func TestDeleteTrigger(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	trigger, err := CreateTrigger(database, automation.ID, TriggerTypeTime, map[string]interface{}{"hour": 9})
	require.NoError(t, err)

	err = DeleteTrigger(database, trigger.ID)
	require.NoError(t, err)

	triggers, err := GetTriggers(database, automation.ID)
	require.NoError(t, err)
	assert.Len(t, triggers, 0)
}

func TestDeleteTriggerNotFound(t *testing.T) {
	database := setupTestDB(t)

	err := DeleteTrigger(database, 999)
	assert.Error(t, err)
}

func TestValidateTriggerType(t *testing.T) {
	testCases := []struct {
		name        string
		triggerType TriggerType
		expectError bool
	}{
		{
			name:        "valid time trigger",
			triggerType: TriggerTypeTime,
			expectError: false,
		},
		{
			name:        "valid sunrise trigger",
			triggerType: TriggerTypeSunrise,
			expectError: false,
		},
		{
			name:        "valid sunset trigger",
			triggerType: TriggerTypeSunset,
			expectError: false,
		},
		{
			name:        "valid presence trigger",
			triggerType: TriggerTypePresence,
			expectError: false,
		},
		{
			name:        "invalid trigger type",
			triggerType: TriggerType("invalid"),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateTriggerType(tc.triggerType)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
