package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/mithilarun/limelight/internal/bridge"
	"github.com/mithilarun/limelight/internal/credentials"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewSetupCommand(logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Initial setup wizard for Hue bridge pairing",
		Long:  "Guides you through the process of connecting to your Hue bridge and storing credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetup(logger)
		},
	}
}

func runSetup(logger *zap.Logger) error {
	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Limelight Setup Wizard ===")
	fmt.Println()

	config, err := credentials.LoadConfig()
	if err != nil {
		return errors.Wrap(err, "loading config")
	}

	if config == nil {
		config = &credentials.Config{}
	}

	fmt.Print("Enter your Hue Bridge IP address: ")
	if config.BridgeIP != "" {
		fmt.Printf("[%s] ", config.BridgeIP)
	}

	bridgeIP, err := reader.ReadString('\n')
	if err != nil {
		return errors.Wrap(err, "reading bridge IP")
	}
	bridgeIP = strings.TrimSpace(bridgeIP)

	if bridgeIP == "" && config.BridgeIP != "" {
		bridgeIP = config.BridgeIP
	}

	if bridgeIP == "" {
		return errors.New("bridge IP is required")
	}
	config.BridgeIP = bridgeIP

	credManager := credentials.NewManager(logger)

	var onePasswordItemName string
	if credManager.IsAvailable(ctx) {
		fmt.Println()
		fmt.Println("1Password CLI detected. API key will be stored securely in 1Password.")
		fmt.Print("Enter 1Password item name [limelight-hue]: ")

		itemName, err := reader.ReadString('\n')
		if err != nil {
			return errors.Wrap(err, "reading item name")
		}
		itemName = strings.TrimSpace(itemName)

		if itemName == "" {
			itemName = "limelight-hue"
		}
		onePasswordItemName = itemName
		config.OnePasswordItemName = itemName
	} else {
		fmt.Println()
		fmt.Println("WARNING: 1Password CLI not available. API key will be stored in plaintext config.")
		fmt.Println("Install 1Password CLI for secure credential storage.")
	}

	fmt.Println()
	fmt.Println("Press the link button on your Hue Bridge now...")
	fmt.Println("Waiting for button press (timeout: 60 seconds)...")
	fmt.Println()

	client := bridge.NewClient(bridgeIP, "", logger)

	authCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	apiKey, err := client.Authenticate(authCtx, "limelight#macOS")
	if err != nil {
		return errors.Wrap(err, "authenticating with bridge")
	}

	if onePasswordItemName != "" {
		if err := credManager.SaveAPIKey(ctx, onePasswordItemName, apiKey); err != nil {
			return errors.Wrap(err, "saving api key to 1password")
		}
		fmt.Println()
		fmt.Printf("API key saved to 1Password item: %s\n", onePasswordItemName)
	}

	if err := credentials.SaveConfig(config); err != nil {
		return errors.Wrap(err, "saving config")
	}

	fmt.Println()
	fmt.Println("Setup complete!")
	fmt.Printf("Bridge IP: %s\n", config.BridgeIP)

	authenticatedClient := bridge.NewClient(bridgeIP, apiKey, logger)
	lights, err := authenticatedClient.GetLights(ctx)
	if err != nil {
		logger.Warn("failed to fetch lights for verification", zap.Error(err))
	} else {
		fmt.Printf("Found %d lights\n", len(lights))
	}

	return nil
}
