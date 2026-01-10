package models

import (
	"database/sql"
	"encoding/json"

	"github.com/cockroachdb/errors"
)

// ActionType represents the type of action
type ActionType string

const (
	ActionTypeLight ActionType = "light"
	ActionTypeScene ActionType = "scene"
	ActionTypeGroup ActionType = "group"
)

// Action represents an action for an automation
type Action struct {
	ID           int64           `json:"id"`
	AutomationID int64           `json:"automation_id"`
	OrderIndex   int             `json:"order_index"`
	Type         ActionType      `json:"type"`
	Config       json.RawMessage `json:"config"`
}

// CreateAction creates a new action
func CreateAction(db *sql.DB, automationID int64, actionType ActionType, config interface{}, orderIndex int) (*Action, error) {
	if err := validateActionType(actionType); err != nil {
		return nil, err
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal action config")
	}

	result, err := db.Exec(
		"INSERT INTO actions (automation_id, type, config, order_index) VALUES (?, ?, ?, ?)",
		automationID, actionType, string(configJSON), orderIndex,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert action")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last insert id")
	}

	return &Action{
		ID:           id,
		AutomationID: automationID,
		OrderIndex:   orderIndex,
		Type:         actionType,
		Config:       configJSON,
	}, nil
}

// GetActions retrieves all actions for an automation, ordered by order_index
func GetActions(db *sql.DB, automationID int64) ([]*Action, error) {
	rows, err := db.Query(
		"SELECT id, automation_id, order_index, type, config FROM actions WHERE automation_id = ? ORDER BY order_index",
		automationID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query actions")
	}
	defer rows.Close()

	var actions []*Action
	for rows.Next() {
		var a Action
		var configStr string
		if err := rows.Scan(&a.ID, &a.AutomationID, &a.OrderIndex, &a.Type, &configStr); err != nil {
			return nil, errors.Wrap(err, "failed to scan action")
		}
		a.Config = json.RawMessage(configStr)
		actions = append(actions, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating actions")
	}

	return actions, nil
}

// DeleteAction deletes an action
func DeleteAction(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM actions WHERE id = ?", id)
	if err != nil {
		return errors.Wrap(err, "failed to delete action")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Newf("action with id %d not found", id)
	}

	return nil
}

// validateActionType validates that the action type is valid
func validateActionType(t ActionType) error {
	switch t {
	case ActionTypeLight, ActionTypeScene, ActionTypeGroup:
		return nil
	default:
		return errors.Newf("invalid action type: %s", t)
	}
}
