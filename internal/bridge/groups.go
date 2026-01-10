package bridge

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
)

type Room struct {
	ID       string `json:"id"`
	IDV1     string `json:"id_v1"`
	Type     string `json:"type"`
	Metadata struct {
		Name      string `json:"name"`
		Archetype string `json:"archetype"`
	} `json:"metadata"`
	Children []struct {
		ResourceID string `json:"rid"`
		Type       string `json:"rtype"`
	} `json:"children"`
}

type RoomsResponse struct {
	Errors []struct {
		Description string `json:"description"`
	} `json:"errors"`
	Data []Room `json:"data"`
}

func (c *Client) GetRooms(ctx context.Context) ([]Room, error) {
	respBody, err := c.doRequest(ctx, "GET", "/resource/room", nil)
	if err != nil {
		return nil, errors.Wrap(err, "getting rooms")
	}

	var roomsResp RoomsResponse
	if err := json.Unmarshal(respBody, &roomsResp); err != nil {
		return nil, errors.Wrap(err, "unmarshaling rooms response")
	}

	if len(roomsResp.Errors) > 0 {
		return nil, errors.Newf("hue api returned errors: %v", roomsResp.Errors)
	}

	return roomsResp.Data, nil
}

type GroupedLight struct {
	ID    string `json:"id"`
	IDV1  string `json:"id_v1"`
	Type  string `json:"type"`
	Owner struct {
		ResourceID string `json:"rid"`
		Type       string `json:"rtype"`
	} `json:"owner"`
	On struct {
		On bool `json:"on"`
	} `json:"on"`
}

type GroupedLightsResponse struct {
	Errors []struct {
		Description string `json:"description"`
	} `json:"errors"`
	Data []GroupedLight `json:"data"`
}

func (c *Client) GetGroupedLights(ctx context.Context) ([]GroupedLight, error) {
	respBody, err := c.doRequest(ctx, "GET", "/resource/grouped_light", nil)
	if err != nil {
		return nil, errors.Wrap(err, "getting grouped lights")
	}

	var groupedLightsResp GroupedLightsResponse
	if err := json.Unmarshal(respBody, &groupedLightsResp); err != nil {
		return nil, errors.Wrap(err, "unmarshaling grouped lights response")
	}

	if len(groupedLightsResp.Errors) > 0 {
		return nil, errors.Newf("hue api returned errors: %v", groupedLightsResp.Errors)
	}

	return groupedLightsResp.Data, nil
}

func (c *Client) SetGroupedLightState(ctx context.Context, groupedLightID string, on bool, brightness *float64) error {
	req := LightUpdateRequest{
		On: &LightOnState{On: on},
	}

	if brightness != nil {
		req.Dimming = &LightDimmingState{
			Brightness: *brightness,
		}
	}

	path := fmt.Sprintf("/resource/grouped_light/%s", groupedLightID)
	_, err := c.doRequest(ctx, "PUT", path, req)
	if err != nil {
		return errors.Wrapf(err, "setting grouped light state for %s", groupedLightID)
	}

	c.logger.Info("grouped light state updated",
		zap.String("grouped_light_id", groupedLightID),
		zap.Bool("on", on),
	)

	return nil
}
