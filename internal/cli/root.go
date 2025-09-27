package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fitboisbot",
	Short: "FitBoisBot - Telegram bot for fitness accountability",
	Long:  `A Telegram bot for fitness accountability and gamification with activity tracking, GG system, and token rewards.`,
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(announceCmd)
	rootCmd.AddCommand(versionCmd)
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

