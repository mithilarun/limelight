package models

import (
	"database/sql"
	"time"

	"github.com/cockroachdb/errors"
)

// Automation represents an automation rule
type Automation struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateAutomation creates a new automation
func CreateAutomation(db *sql.DB, name, description string) (*Automation, error) {
	if name == "" {
		return nil, errors.New("automation name cannot be empty")
	}

	result, err := db.Exec(
		"INSERT INTO automations (name, description) VALUES (?, ?)",
		name, description,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert automation")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last insert id")
	}

	return GetAutomation(db, id)
}

// GetAutomation retrieves an automation by ID
func GetAutomation(db *sql.DB, id int64) (*Automation, error) {
	var a Automation
	err := db.QueryRow(
		"SELECT id, name, description, enabled, created_at, updated_at FROM automations WHERE id = ?",
		id,
	).Scan(&a.ID, &a.Name, &a.Description, &a.Enabled, &a.CreatedAt, &a.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.Newf("automation with id %d not found", id)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to query automation")
	}

	return &a, nil
}

// ListAutomations retrieves all automations
func ListAutomations(db *sql.DB) ([]*Automation, error) {
	rows, err := db.Query(
		"SELECT id, name, description, enabled, created_at, updated_at FROM automations ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query automations")
	}
	defer rows.Close()

	var automations []*Automation
	for rows.Next() {
		var a Automation
		if err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.Enabled, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan automation")
		}
		automations = append(automations, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating automations")
	}

	return automations, nil
}

// ListEnabledAutomations retrieves all enabled automations
func ListEnabledAutomations(db *sql.DB) ([]*Automation, error) {
	rows, err := db.Query(
		"SELECT id, name, description, enabled, created_at, updated_at FROM automations WHERE enabled = 1 ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query enabled automations")
	}
	defer rows.Close()

	var automations []*Automation
	for rows.Next() {
		var a Automation
		if err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.Enabled, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan automation")
		}
		automations = append(automations, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating enabled automations")
	}

	return automations, nil
}

// UpdateAutomation updates an automation's name and description
func UpdateAutomation(db *sql.DB, id int64, name, description string) error {
	if name == "" {
		return errors.New("automation name cannot be empty")
	}

	result, err := db.Exec(
		"UPDATE automations SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		name, description, id,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update automation")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Newf("automation with id %d not found", id)
	}

	return nil
}

// SetEnabled enables or disables an automation
func SetEnabled(db *sql.DB, id int64, enabled bool) error {
	result, err := db.Exec(
		"UPDATE automations SET enabled = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		enabled, id,
	)
	if err != nil {
		return errors.Wrap(err, "failed to set automation enabled status")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Newf("automation with id %d not found", id)
	}

	return nil
}

// DeleteAutomation deletes an automation and all associated triggers, conditions, and actions.
// Due to CASCADE DELETE foreign key constraints, all related triggers, conditions,
// and actions are automatically deleted when the automation is deleted.
func DeleteAutomation(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM automations WHERE id = ?", id)
	if err != nil {
		return errors.Wrap(err, "failed to delete automation")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Newf("automation with id %d not found", id)
	}

	return nil
}
