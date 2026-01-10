package bridge

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
)

type Scene struct {
	ID       string `json:"id"`
	IDV1     string `json:"id_v1"`
	Type     string `json:"type"`
	Metadata struct {
		Name  string `json:"name"`
		Image *struct {
			ResourceID string `json:"rid"`
			Type       string `json:"rtype"`
		} `json:"image,omitempty"`
	} `json:"metadata"`
	Group struct {
		ResourceID string `json:"rid"`
		Type       string `json:"rtype"`
	} `json:"group"`
	Actions []struct {
		Target struct {
			ResourceID string `json:"rid"`
			Type       string `json:"rtype"`
		} `json:"target"`
		Action struct {
			On      *LightOnState      `json:"on,omitempty"`
			Dimming *LightDimmingState `json:"dimming,omitempty"`
		} `json:"action"`
	} `json:"actions"`
}

type ScenesResponse struct {
	Errors []struct {
		Description string `json:"description"`
	} `json:"errors"`
	Data []Scene `json:"data"`
}

type SceneRecallRequest struct {
	Recall struct {
		Action string `json:"action"`
	} `json:"recall"`
}

func (c *Client) GetScenes(ctx context.Context) ([]Scene, error) {
	respBody, err := c.doRequest(ctx, "GET", "/resource/scene", nil)
	if err != nil {
		return nil, errors.Wrap(err, "getting scenes")
	}

	var scenesResp ScenesResponse
	if err := json.Unmarshal(respBody, &scenesResp); err != nil {
		return nil, errors.Wrap(err, "unmarshaling scenes response")
	}

	if len(scenesResp.Errors) > 0 {
		return nil, errors.Newf("hue api returned errors: %v", scenesResp.Errors)
	}

	return scenesResp.Data, nil
}

func (c *Client) ActivateScene(ctx context.Context, sceneID string) error {
	req := SceneRecallRequest{}
	req.Recall.Action = "active"

	path := fmt.Sprintf("/resource/scene/%s", sceneID)
	_, err := c.doRequest(ctx, "PUT", path, req)
	if err != nil {
		return errors.Wrapf(err, "activating scene %s", sceneID)
	}

	c.logger.Info("scene activated",
		zap.String("scene_id", sceneID),
	)

	return nil
}
