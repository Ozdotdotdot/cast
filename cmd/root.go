package cmd

import (
	"fmt"
	"os"

	"github.com/Ozdotdotdot/cast/assets"
	"github.com/Ozdotdotdot/cast/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cast [device] <message>",
	Short: "Send text-to-speech messages to Google Home devices",
	Long: `CAST-CLI is a friendly CLI tool for sending text-to-speech messages
to Google Home and Chromecast devices on your network.

Examples:
  cast "Hello world"              # Speak to default device
  cast bedroom "Goodnight"        # Speak to specific device
  cast hq,bedroom "Dinner time"   # Speak to multiple devices
  cast everywhere "Attention"     # Speak to a group`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			// First run - no config exists
			if os.IsNotExist(err) {
				return runFirstTimeSetup(args)
			}
			return fmt.Errorf("failed to load config: %w", err)
		}

		var device, message string
		if len(args) == 1 {
			device = cfg.DefaultDevice
			message = args[0]
			if device == "" {
				return fmt.Errorf("no default device configured. Run 'cast discover' or specify a device")
			}
		} else {
			device = args[0]
			message = args[1]
		}

		return speak(cfg, device, message)
	},
}

func runFirstTimeSetup(args []string) error {
	fmt.Print(assets.Banner)
	fmt.Println(assets.Welcome)
	fmt.Println()

	// Run discovery
	devices, err := runDiscoveryFlow()
	if err != nil {
		return err
	}

	if len(devices) == 0 {
		fmt.Println("No devices found. Make sure your Google Home devices are on the same network.")
		fmt.Println("You can run 'cast discover' later to try again.")
		return nil
	}

	// Save config
	cfg, err := saveDiscoveredDevices(devices)
	if err != nil {
		return err
	}

	// User declined to save
	if cfg == nil {
		return nil
	}

	// If user provided a message, speak it
	if len(args) >= 1 {
		var device, message string
		if len(args) == 1 {
			device = cfg.DefaultDevice
			message = args[0]
		} else {
			device = args[0]
			message = args[1]
		}
		if device != "" && message != "" {
			fmt.Println()
			return speak(cfg, device, message)
		}
	}

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(discoverCmd)
	rootCmd.AddCommand(devicesCmd)
	rootCmd.AddCommand(groupsCmd)
	rootCmd.AddCommand(configCmd)
}
