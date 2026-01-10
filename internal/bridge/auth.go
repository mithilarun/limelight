package bridge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
)

type authRequest struct {
	DeviceType        string `json:"devicetype"`
	GenerateClientKey bool   `json:"generateclientkey"`
}

type authResponseItem struct {
	Success *struct {
		Username string `json:"username"`
	} `json:"success,omitempty"`
	Error *struct {
		Type        int    `json:"type"`
		Address     string `json:"address"`
		Description string `json:"description"`
	} `json:"error,omitempty"`
}

const (
	authErrorLinkButtonNotPressed = 101
	authRetryInterval             = 2 * time.Second
	authMaxRetries                = 30
)

func (c *Client) Authenticate(ctx context.Context, appName string) (string, error) {
	url := fmt.Sprintf("https://%s/api", c.bridgeIP)

	reqBody := authRequest{
		DeviceType:        appName,
		GenerateClientKey: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", errors.Wrap(err, "marshaling auth request")
	}

	var lastError error
	for i := 0; i < authMaxRetries; i++ {
		select {
		case <-ctx.Done():
			return "", errors.Wrap(ctx.Err(), "authentication cancelled")
		default:
		}

		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonData))
		if err != nil {
			return "", errors.Wrap(err, "creating auth request")
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return "", errors.Wrap(err, "executing auth request")
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return "", errors.Wrap(err, "reading auth response")
		}

		var authResp []authResponseItem
		if err := json.Unmarshal(respBody, &authResp); err != nil {
			return "", errors.Wrap(err, "unmarshaling auth response")
		}

		if len(authResp) == 0 {
			return "", errors.New("empty auth response")
		}

		item := authResp[0]

		if item.Success != nil {
			c.logger.Info("authentication successful",
				zap.String("username", item.Success.Username),
			)
			return item.Success.Username, nil
		}

		if item.Error != nil {
			if item.Error.Type == authErrorLinkButtonNotPressed {
				c.logger.Debug("waiting for link button press",
					zap.Int("attempt", i+1),
					zap.Int("max_retries", authMaxRetries),
				)
				lastError = errors.Newf("link button not pressed: %s", item.Error.Description)
				time.Sleep(authRetryInterval)
				continue
			}
			return "", errors.Newf("hue auth error: type=%d, description=%s", item.Error.Type, item.Error.Description)
		}
	}

	if lastError != nil {
		return "", errors.Wrap(lastError, "authentication timed out")
	}
	return "", errors.New("authentication failed")
}
