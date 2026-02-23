package cmd

import (
	"fmt"

	"github.com/Ozdotdotdot/cast/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show config file location",
	Long:  `Display the path to the CAST-CLI configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Config file:", config.Path())
	},
}
