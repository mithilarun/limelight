package commands

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewScenesCommand(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scenes",
		Short: "Manage Hue scenes",
		Long:  "List and activate Hue scenes",
	}

	cmd.AddCommand(newListScenesCommand(logger))
	cmd.AddCommand(newActivateSceneCommand(logger))

	return cmd
}

func newListScenesCommand(logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all scenes",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := getAuthenticatedClient(ctx, logger)
			if err != nil {
				return err
			}

			scenes, err := client.GetScenes(ctx)
			if err != nil {
				return errors.Wrap(err, "getting scenes")
			}

			fmt.Printf("Found %d scenes:\n\n", len(scenes))
			for _, scene := range scenes {
				fmt.Printf("  %s\n", scene.Metadata.Name)
				fmt.Printf("    ID: %s\n", scene.ID)
				fmt.Printf("    Group: %s\n", scene.Group.ResourceID)
				fmt.Println()
			}

			return nil
		},
	}
}

func newActivateSceneCommand(logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "activate <scene-id>",
		Short: "Activate a scene",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := getAuthenticatedClient(ctx, logger)
			if err != nil {
				return err
			}

			sceneID := args[0]

			if err := client.ActivateScene(ctx, sceneID); err != nil {
				return errors.Wrap(err, "activating scene")
			}

			fmt.Printf("Scene %s activated successfully\n", sceneID)
			return nil
		},
	}
}
