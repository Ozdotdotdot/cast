package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Ozdotdotdot/cast/internal/config"
	"github.com/spf13/cobra"
)

var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List configured groups",
	Long:  `Display all configured device groups.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("no config found. Run 'cast discover' first")
		}

		if len(cfg.Groups) == 0 {
			fmt.Println("No groups configured.")
			fmt.Println("Add groups to your config file at:", config.Path())
			return nil
		}

		// Sort group names for consistent output
		var names []string
		for name := range cfg.Groups {
			names = append(names, name)
		}
		sort.Strings(names)

		fmt.Println("Configured groups:")
		for _, name := range names {
			g := cfg.Groups[name]
			fmt.Printf("  %s: %s\n", name, strings.Join(g.Devices, ", "))
		}

		return nil
	},
}
