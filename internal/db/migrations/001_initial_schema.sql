-- Create automations table
CREATE TABLE automations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    enabled BOOLEAN NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create triggers table
CREATE TABLE triggers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    automation_id INTEGER NOT NULL,
    type TEXT NOT NULL,
    config TEXT NOT NULL,
    FOREIGN KEY (automation_id) REFERENCES automations(id) ON DELETE CASCADE
);

-- Create conditions table
CREATE TABLE conditions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    automation_id INTEGER NOT NULL,
    type TEXT NOT NULL,
    config TEXT NOT NULL,
    FOREIGN KEY (automation_id) REFERENCES automations(id) ON DELETE CASCADE
);

-- Create actions table
CREATE TABLE actions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    automation_id INTEGER NOT NULL,
    order_index INTEGER NOT NULL DEFAULT 0,
    type TEXT NOT NULL,
    config TEXT NOT NULL,
    FOREIGN KEY (automation_id) REFERENCES automations(id) ON DELETE CASCADE
);

-- Create config table
CREATE TABLE config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indices for foreign keys
CREATE INDEX idx_triggers_automation_id ON triggers(automation_id);
CREATE INDEX idx_conditions_automation_id ON conditions(automation_id);
CREATE INDEX idx_actions_automation_id ON actions(automation_id);

-- Create index on automations.enabled for quick lookup of active automations
CREATE INDEX idx_automations_enabled ON automations(enabled);

-- Create index on actions.order_index for proper action ordering
CREATE INDEX idx_actions_order ON actions(automation_id, order_index);
