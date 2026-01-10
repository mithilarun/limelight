package bridge

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
)

type Light struct {
	ID       string `json:"id"`
	IDV1     string `json:"id_v1"`
	Type     string `json:"type"`
	Metadata struct {
		Name      string `json:"name"`
		Archetype string `json:"archetype"`
	} `json:"metadata"`
	On struct {
		On bool `json:"on"`
	} `json:"on"`
	Dimming *struct {
		Brightness float64 `json:"brightness"`
	} `json:"dimming,omitempty"`
	ColorTemperature *struct {
		Mirek int `json:"mirek"`
	} `json:"color_temperature,omitempty"`
	Color *struct {
		XY struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"xy"`
	} `json:"color,omitempty"`
}

type LightsResponse struct {
	Errors []struct {
		Description string `json:"description"`
	} `json:"errors"`
	Data []Light `json:"data"`
}

type LightUpdateRequest struct {
	On      *LightOnState      `json:"on,omitempty"`
	Dimming *LightDimmingState `json:"dimming,omitempty"`
}

type LightOnState struct {
	On bool `json:"on"`
}

type LightDimmingState struct {
	Brightness float64 `json:"brightness"`
}

func (c *Client) GetLights(ctx context.Context) ([]Light, error) {
	respBody, err := c.doRequest(ctx, "GET", "/resource/light", nil)
	if err != nil {
		return nil, errors.Wrap(err, "getting lights")
	}

	var lightsResp LightsResponse
	if err := json.Unmarshal(respBody, &lightsResp); err != nil {
		return nil, errors.Wrap(err, "unmarshaling lights response")
	}

	if len(lightsResp.Errors) > 0 {
		return nil, errors.Newf("hue api returned errors: %v", lightsResp.Errors)
	}

	return lightsResp.Data, nil
}

func (c *Client) SetLightState(ctx context.Context, lightID string, on bool, brightness *float64) error {
	req := LightUpdateRequest{
		On: &LightOnState{On: on},
	}

	if brightness != nil {
		req.Dimming = &LightDimmingState{
			Brightness: *brightness,
		}
	}

	path := fmt.Sprintf("/resource/light/%s", lightID)
	_, err := c.doRequest(ctx, "PUT", path, req)
	if err != nil {
		return errors.Wrapf(err, "setting light state for light %s", lightID)
	}

	c.logger.Info("light state updated",
		zap.String("light_id", lightID),
		zap.Bool("on", on),
	)

	return nil
}
