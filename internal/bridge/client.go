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

const (
	apiVersion = "v2"
	httpTimeout = 10 * time.Second
)

type Client struct {
	bridgeIP   string
	apiKey     string
	httpClient *http.Client
	logger     *zap.Logger
}

func NewClient(bridgeIP, apiKey string, logger *zap.Logger) *Client {
	return &Client{
		bridgeIP: bridgeIP,
		apiKey:   apiKey,
		httpClient: &http.Client{
			Timeout: httpTimeout,
		},
		logger: logger,
	}
}

func (c *Client) baseURL() string {
	return fmt.Sprintf("https://%s/clip/%s", c.bridgeIP, apiVersion)
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	url := c.baseURL() + path

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(err, "marshaling request body")
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "creating http request")
	}

	if c.apiKey != "" {
		req.Header.Set("hue-application-key", c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	c.logger.Debug("hue api request",
		zap.String("method", method),
		zap.String("url", url),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "executing http request")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.Newf("hue api error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
