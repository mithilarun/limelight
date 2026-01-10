package main

import (
	"fmt"
	"os"

	"github.com/mithilarun/limelight/cmd/limelight/commands"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	rootCmd := &cobra.Command{
		Use:   "limelight",
		Short: "Philips Hue automation tool",
		Long:  "A CLI tool for proactive automation of Philips Hue lights and scenes",
	}

	rootCmd.AddCommand(commands.NewSetupCommand(logger))
	rootCmd.AddCommand(commands.NewLightsCommand(logger))
	rootCmd.AddCommand(commands.NewScenesCommand(logger))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
