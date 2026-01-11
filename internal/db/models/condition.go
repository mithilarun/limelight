package models

import (
	"database/sql"
	"encoding/json"

	"github.com/cockroachdb/errors"
)

// ConditionType represents the type of condition
type ConditionType string

const (
	ConditionTypeWeekday   ConditionType = "weekday"
	ConditionTypeWeekend   ConditionType = "weekend"
	ConditionTypeDayOfWeek ConditionType = "day_of_week"
	ConditionTypeDateRange ConditionType = "date_range"
)

// Condition represents a condition for an automation
type Condition struct {
	ID           int64           `json:"id"`
	AutomationID int64           `json:"automation_id"`
	Type         ConditionType   `json:"type"`
	Config       json.RawMessage `json:"config"`
}

// CreateCondition creates a new condition
func CreateCondition(db *sql.DB, automationID int64, conditionType ConditionType, config interface{}) (*Condition, error) {
	if err := validateConditionType(conditionType); err != nil {
		return nil, err
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal condition config")
	}

	result, err := db.Exec(
		"INSERT INTO conditions (automation_id, type, config) VALUES (?, ?, ?)",
		automationID, conditionType, string(configJSON),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert condition")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last insert id")
	}

	return &Condition{
		ID:           id,
		AutomationID: automationID,
		Type:         conditionType,
		Config:       configJSON,
	}, nil
}

// GetConditions retrieves all conditions for an automation
func GetConditions(db *sql.DB, automationID int64) ([]*Condition, error) {
	rows, err := db.Query(
		"SELECT id, automation_id, type, config FROM conditions WHERE automation_id = ?",
		automationID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query conditions")
	}
	defer rows.Close()

	var conditions []*Condition
	for rows.Next() {
		var c Condition
		var configStr string
		if err := rows.Scan(&c.ID, &c.AutomationID, &c.Type, &configStr); err != nil {
			return nil, errors.Wrap(err, "failed to scan condition")
		}
		c.Config = json.RawMessage(configStr)
		conditions = append(conditions, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating conditions")
	}

	return conditions, nil
}

// DeleteCondition deletes a condition
func DeleteCondition(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM conditions WHERE id = ?", id)
	if err != nil {
		return errors.Wrap(err, "failed to delete condition")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Newf("condition with id %d not found", id)
	}

	return nil
}

// validateConditionType validates that the condition type is valid
func validateConditionType(t ConditionType) error {
	switch t {
	case ConditionTypeWeekday, ConditionTypeWeekend, ConditionTypeDayOfWeek, ConditionTypeDateRange:
		return nil
	default:
		return errors.Newf("invalid condition type: %s", t)
	}
}
