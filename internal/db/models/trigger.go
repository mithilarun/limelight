package models

import (
	"database/sql"
	"encoding/json"

	"github.com/cockroachdb/errors"
)

// TriggerType represents the type of trigger
type TriggerType string

const (
	TriggerTypeTime     TriggerType = "time"
	TriggerTypeSunrise  TriggerType = "sunrise"
	TriggerTypeSunset   TriggerType = "sunset"
	TriggerTypePresence TriggerType = "presence"
)

// Trigger represents a trigger for an automation
type Trigger struct {
	ID           int64           `json:"id"`
	AutomationID int64           `json:"automation_id"`
	Type         TriggerType     `json:"type"`
	Config       json.RawMessage `json:"config"`
}

// CreateTrigger creates a new trigger
func CreateTrigger(db *sql.DB, automationID int64, triggerType TriggerType, config interface{}) (*Trigger, error) {
	if err := validateTriggerType(triggerType); err != nil {
		return nil, err
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal trigger config")
	}

	result, err := db.Exec(
		"INSERT INTO triggers (automation_id, type, config) VALUES (?, ?, ?)",
		automationID, triggerType, string(configJSON),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert trigger")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last insert id")
	}

	return &Trigger{
		ID:           id,
		AutomationID: automationID,
		Type:         triggerType,
		Config:       configJSON,
	}, nil
}

// GetTriggers retrieves all triggers for an automation
func GetTriggers(db *sql.DB, automationID int64) ([]*Trigger, error) {
	rows, err := db.Query(
		"SELECT id, automation_id, type, config FROM triggers WHERE automation_id = ?",
		automationID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query triggers")
	}
	defer rows.Close()

	var triggers []*Trigger
	for rows.Next() {
		var t Trigger
		var configStr string
		if err := rows.Scan(&t.ID, &t.AutomationID, &t.Type, &configStr); err != nil {
			return nil, errors.Wrap(err, "failed to scan trigger")
		}
		t.Config = json.RawMessage(configStr)
		triggers = append(triggers, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating triggers")
	}

	return triggers, nil
}

// DeleteTrigger deletes a trigger
func DeleteTrigger(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM triggers WHERE id = ?", id)
	if err != nil {
		return errors.Wrap(err, "failed to delete trigger")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Newf("trigger with id %d not found", id)
	}

	return nil
}

// validateTriggerType validates that the trigger type is valid
func validateTriggerType(t TriggerType) error {
	switch t {
	case TriggerTypeTime, TriggerTypeSunrise, TriggerTypeSunset, TriggerTypePresence:
		return nil
	default:
		return errors.Newf("invalid trigger type: %s", t)
	}
}
