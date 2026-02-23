package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Ozdotdotdot/cast/internal/config"
	"github.com/spf13/cobra"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "List configured devices",
	Long:  `Display all configured Google Home devices and their aliases.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("no config found. Run 'cast discover' first")
		}

		if len(cfg.Devices) == 0 {
			fmt.Println("No devices configured. Run 'cast discover' to find devices.")
			return nil
		}

		// Sort device names for consistent output
		var names []string
		for name := range cfg.Devices {
			names = append(names, name)
		}
		sort.Strings(names)

		fmt.Println("Configured devices:")
		for _, name := range names {
			d := cfg.Devices[name]
			defaultMarker := ""
			if name == cfg.DefaultDevice {
				defaultMarker = " (default)"
			}

			aliases := ""
			if len(d.Aliases) > 0 {
				aliases = fmt.Sprintf(" [aliases: %s]", strings.Join(d.Aliases, ", "))
			}

			fmt.Printf("  %s%s%s\n", name, defaultMarker, aliases)
			fmt.Printf("    %s (%s:%d)\n", d.OriginalName, d.IP, d.Port)
		}

		return nil
	},
}
