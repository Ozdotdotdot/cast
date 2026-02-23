package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Ozdotdotdot/cast/internal/config"
	"github.com/Ozdotdotdot/cast/internal/discovery"
	"github.com/spf13/cobra"
)

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Scan for Google Home devices on your network",
	Long:  `Scans your local network for Google Home and Chromecast devices using mDNS.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		devices, err := runDiscoveryFlow()
		if err != nil {
			return err
		}

		if len(devices) == 0 {
			fmt.Println("No devices found. Make sure your Google Home devices are on the same network.")
			return nil
		}

		_, err = saveDiscoveredDevices(devices)
		return err
	},
}

func runDiscoveryFlow() ([]discovery.DiscoveredDevice, error) {
	fmt.Println("Scanning for devices...")
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	ch, err := discovery.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}

	var devices []discovery.DiscoveredDevice
	for d := range ch {
		devices = append(devices, d)
		fmt.Printf("  Found: %s (%s)\n", d.Name, d.IP)
	}

	if len(devices) == 0 {
		return nil, nil
	}

	fmt.Println()
	fmt.Printf("Found %d device(s).\n", len(devices))
	fmt.Println()

	return devices, nil
}

func saveDiscoveredDevices(devices []discovery.DiscoveredDevice) (*config.Config, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Add these devices to config? [Y/n]: ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "" && response != "y" && response != "yes" {
		fmt.Println("Devices not saved.")
		return nil, nil
	}

	// Load existing config or create new
	cfg, err := config.Load()
	if err != nil {
		cfg = config.New()
	}

	// Add devices
	for _, d := range devices {
		normalizedName := normalizeName(d.Name)
		cfg.Devices[normalizedName] = config.Device{
			IP:           d.IP,
			Port:         d.Port,
			OriginalName: d.Name,
			Aliases:      []string{},
		}
	}

	// Set default device if not set
	if cfg.DefaultDevice == "" && len(devices) > 0 {
		defaultName := normalizeName(devices[0].Name)
		fmt.Printf("Set default device? [%s]: ", defaultName)
		response, _ = reader.ReadString('\n')
		response = strings.TrimSpace(response)

		if response == "" {
			cfg.DefaultDevice = defaultName
		} else {
			cfg.DefaultDevice = response
		}
	}

	// Save config
	if err := cfg.Save(); err != nil {
		return nil, fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Devices saved to %s\n", config.Path())

	return cfg, nil
}

// normalizeName converts a device name to a CLI-friendly identifier
// "Headquarters Speaker" -> "headquarters-speaker"
func normalizeName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace spaces and underscores with hyphens
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")

	// Remove any characters that aren't alphanumeric or hyphens
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	name = reg.ReplaceAllString(name, "")

	// Collapse multiple hyphens
	reg = regexp.MustCompile(`-+`)
	name = reg.ReplaceAllString(name, "-")

	// Trim leading/trailing hyphens
	name = strings.Trim(name, "-")

	return name
}
