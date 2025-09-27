package cli

import (
	"fmt"

	"github.com/Alcamech/FitBoisBot/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display the current version of FitBoisBot along with build information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.GetDetailedVersionInfo())
	},
}