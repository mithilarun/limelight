package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCondition(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	config := map[string]interface{}{
		"days": []int{1, 2, 3, 4, 5},
	}

	condition, err := CreateCondition(database, automation.ID, ConditionTypeDayOfWeek, config)
	require.NoError(t, err)
	require.NotNil(t, condition)

	assert.Greater(t, condition.ID, int64(0))
	assert.Equal(t, automation.ID, condition.AutomationID)
	assert.Equal(t, ConditionTypeDayOfWeek, condition.Type)

	var parsed map[string]interface{}
	err = json.Unmarshal(condition.Config, &parsed)
	require.NoError(t, err)
	assert.NotNil(t, parsed["days"])
}

func TestCreateConditionInvalidType(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	_, err = CreateCondition(database, automation.ID, ConditionType("invalid"), map[string]interface{}{})
	assert.Error(t, err)
}

func TestGetConditions(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	_, err = CreateCondition(database, automation.ID, ConditionTypeWeekday, map[string]interface{}{})
	require.NoError(t, err)

	_, err = CreateCondition(database, automation.ID, ConditionTypeWeekend, map[string]interface{}{})
	require.NoError(t, err)

	conditions, err := GetConditions(database, automation.ID)
	require.NoError(t, err)
	assert.Len(t, conditions, 2)
}

func TestGetConditionsEmpty(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	conditions, err := GetConditions(database, automation.ID)
	require.NoError(t, err)
	assert.Len(t, conditions, 0)
}

func TestDeleteCondition(t *testing.T) {
	database := setupTestDB(t)

	automation, err := CreateAutomation(database, "Test", "Test automation")
	require.NoError(t, err)

	condition, err := CreateCondition(database, automation.ID, ConditionTypeWeekday, map[string]interface{}{})
	require.NoError(t, err)

	err = DeleteCondition(database, condition.ID)
	require.NoError(t, err)

	conditions, err := GetConditions(database, automation.ID)
	require.NoError(t, err)
	assert.Len(t, conditions, 0)
}

func TestDeleteConditionNotFound(t *testing.T) {
	database := setupTestDB(t)

	err := DeleteCondition(database, 999)
	assert.Error(t, err)
}

func TestValidateConditionType(t *testing.T) {
	testCases := []struct {
		name          string
		conditionType ConditionType
		expectError   bool
	}{
		{
			name:          "valid weekday condition",
			conditionType: ConditionTypeWeekday,
			expectError:   false,
		},
		{
			name:          "valid weekend condition",
			conditionType: ConditionTypeWeekend,
			expectError:   false,
		},
		{
			name:          "valid day_of_week condition",
			conditionType: ConditionTypeDayOfWeek,
			expectError:   false,
		},
		{
			name:          "valid date_range condition",
			conditionType: ConditionTypeDateRange,
			expectError:   false,
		},
		{
			name:          "invalid condition type",
			conditionType: ConditionType("invalid"),
			expectError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateConditionType(tc.conditionType)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
