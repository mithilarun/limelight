package credentials

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
)

const (
	onePasswordCLI = "op"
)

type Manager struct {
	logger *zap.Logger
}

func NewManager(logger *zap.Logger) *Manager {
	return &Manager{
		logger: logger,
	}
}

type opItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Fields []struct {
		ID    string `json:"id"`
		Type  string `json:"type"`
		Label string `json:"label"`
		Value string `json:"value"`
	} `json:"fields"`
}

func (m *Manager) GetAPIKey(ctx context.Context, itemName string) (string, error) {
	cmd := exec.CommandContext(ctx, onePasswordCLI, "item", "get", itemName, "--format", "json")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	m.logger.Debug("executing 1password cli command",
		zap.String("item", itemName),
	)

	if err := cmd.Run(); err != nil {
		return "", errors.Wrapf(err, "executing op cli command: stderr=%s", stderr.String())
	}

	var item opItem
	if err := json.Unmarshal(stdout.Bytes(), &item); err != nil {
		return "", errors.Wrap(err, "unmarshaling op item response")
	}

	for _, field := range item.Fields {
		if field.Label == "api_key" || field.Label == "apikey" || field.Type == "CONCEALED" {
			if field.Value != "" {
				m.logger.Info("api key retrieved from 1password",
					zap.String("item", itemName),
				)
				return field.Value, nil
			}
		}
	}

	return "", errors.Newf("api key field not found in 1password item %s", itemName)
}

func (m *Manager) SaveAPIKey(ctx context.Context, itemName, apiKey string) error {
	getCmd := exec.CommandContext(ctx, onePasswordCLI, "item", "get", itemName)
	if err := getCmd.Run(); err == nil {
		editCmd := exec.CommandContext(ctx, onePasswordCLI, "item", "edit", itemName, fmt.Sprintf("api_key[concealed]=%s", apiKey))

		var stderr bytes.Buffer
		editCmd.Stderr = &stderr

		if err := editCmd.Run(); err != nil {
			return errors.Wrapf(err, "editing 1password item: stderr=%s", stderr.String())
		}

		m.logger.Info("api key updated in 1password",
			zap.String("item", itemName),
		)
		return nil
	}

	createCmd := exec.CommandContext(ctx, onePasswordCLI, "item", "create",
		"--category", "api_credential",
		"--title", itemName,
		fmt.Sprintf("api_key[concealed]=%s", apiKey),
	)

	var stderr bytes.Buffer
	createCmd.Stderr = &stderr

	if err := createCmd.Run(); err != nil {
		return errors.Wrapf(err, "creating 1password item: stderr=%s", stderr.String())
	}

	m.logger.Info("api key saved to 1password",
		zap.String("item", itemName),
	)

	return nil
}

func (m *Manager) IsAvailable(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "which", onePasswordCLI)
	if err := cmd.Run(); err != nil {
		return false
	}

	cmd = exec.CommandContext(ctx, onePasswordCLI, "account", "list")
	if err := cmd.Run(); err != nil {
		m.logger.Warn("1password cli is installed but not signed in")
		return false
	}

	return true
}
