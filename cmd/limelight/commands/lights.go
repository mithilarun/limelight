package commands

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/mithilarun/limelight/internal/bridge"
	"github.com/mithilarun/limelight/internal/credentials"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewLightsCommand(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lights",
		Short: "Manage Hue lights",
		Long:  "List and control Hue lights",
	}

	cmd.AddCommand(newListLightsCommand(logger))
	cmd.AddCommand(newSetLightCommand(logger))

	return cmd
}

func newListLightsCommand(logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all lights",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := getAuthenticatedClient(ctx, logger)
			if err != nil {
				return err
			}

			lights, err := client.GetLights(ctx)
			if err != nil {
				return errors.Wrap(err, "getting lights")
			}

			fmt.Printf("Found %d lights:\n\n", len(lights))
			for _, light := range lights {
				status := "off"
				if light.On.On {
					status = "on"
				}

				brightness := ""
				if light.Dimming != nil {
					brightness = fmt.Sprintf(" (%.0f%%)", light.Dimming.Brightness)
				}

				fmt.Printf("  %s\n", light.Metadata.Name)
				fmt.Printf("    ID: %s\n", light.ID)
				fmt.Printf("    Status: %s%s\n", status, brightness)
				fmt.Printf("    Type: %s\n", light.Metadata.Archetype)
				fmt.Println()
			}

			return nil
		},
	}
}

func newSetLightCommand(logger *zap.Logger) *cobra.Command {
	var (
		on         bool
		off        bool
		brightness float64
	)

	cmd := &cobra.Command{
		Use:   "set <light-id>",
		Short: "Set light state",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := getAuthenticatedClient(ctx, logger)
			if err != nil {
				return err
			}

			lightID := args[0]

			targetOn := on
			if off {
				targetOn = false
			}

			var brightnessPtr *float64
			if brightness > 0 {
				brightnessPtr = &brightness
			}

			if err := client.SetLightState(ctx, lightID, targetOn, brightnessPtr); err != nil {
				return errors.Wrap(err, "setting light state")
			}

			fmt.Printf("Light %s updated successfully\n", lightID)
			return nil
		},
	}

	cmd.Flags().BoolVar(&on, "on", false, "Turn light on")
	cmd.Flags().BoolVar(&off, "off", false, "Turn light off")
	cmd.Flags().Float64Var(&brightness, "brightness", 0, "Set brightness (0-100)")

	return cmd
}

func getAuthenticatedClient(ctx context.Context, logger *zap.Logger) (*bridge.Client, error) {
	config, err := credentials.LoadConfig()
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	if config == nil {
		return nil, errors.New("no configuration found, run 'limelight setup' first")
	}

	var apiKey string

	if config.OnePasswordItemName != "" {
		credManager := credentials.NewManager(logger)
		if !credManager.IsAvailable(ctx) {
			return nil, errors.New("1password cli not available")
		}

		apiKey, err = credManager.GetAPIKey(ctx, config.OnePasswordItemName)
		if err != nil {
			return nil, errors.Wrap(err, "getting api key from 1password")
		}
	} else {
		return nil, errors.New("no credential storage configured")
	}

	return bridge.NewClient(config.BridgeIP, apiKey, logger), nil
}
