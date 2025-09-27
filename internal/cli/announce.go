package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Alcamech/FitBoisBot/config"
	"github.com/Alcamech/FitBoisBot/internal/bot"
	"github.com/Alcamech/FitBoisBot/internal/database"
	"github.com/spf13/cobra"
)

var (
	announceFile    string
	announceGroup   int64
	announcePreview bool
)

var announceCmd = &cobra.Command{
	Use:   "announce [message]",
	Short: "Send announcement to groups",
	Long: `Send an announcement message to groups. Can send to all groups or a specific group.
	
Examples:
  fitboisbot announce "Bot updated with new features!"
  fitboisbot announce --file announcement.md
  fitboisbot announce --group -123456789 "Group-specific message"
  fitboisbot announce --preview --file update.md`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Must have either message arg or --file flag
		if len(args) == 0 && announceFile == "" {
			return fmt.Errorf("must provide either a message argument or --file flag")
		}
		if len(args) > 0 && announceFile != "" {
			return fmt.Errorf("cannot use both message argument and --file flag")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var message string
		var err error

		// Get message content
		if announceFile != "" {
			content, err := os.ReadFile(announceFile)
			if err != nil {
				slog.Error("Failed to read file", "error", err, "file", announceFile)
				os.Exit(1)
			}
			message = string(content)
		} else {
			message = args[0]
		}

		// Preview mode - just show what would be sent
		if announcePreview {
			fmt.Printf("Preview mode - would send the following message:\n\n%s\n", message)
			if announceGroup != 0 {
				fmt.Printf("\nTarget: Group %d\n", announceGroup)
			} else {
				fmt.Printf("\nTarget: All groups\n")
			}
			return
		}

		// Initialize services
		config.InitConfig()
		database.InitDB()

		service, err := bot.NewBotService()
		if err != nil {
			slog.Error("Failed to initialize bot service", "error", err)
			os.Exit(1)
		}

		// Send to specific group
		if announceGroup != 0 {
			err := service.SendAnnouncement(announceGroup, message)
			if err != nil {
				slog.Error("Failed to send announcement", "error", err, "group_id", announceGroup)
				os.Exit(1)
			}
			fmt.Printf("Announcement sent to group %d\n", announceGroup)
			return
		}

		// Send to all groups
		groups, err := service.GetAllGroups()
		if err != nil {
			slog.Error("Failed to fetch groups", "error", err)
			os.Exit(1)
		}

		slog.Info("Sending announcement to groups", "count", len(groups))

		successCount := 0
		for _, group := range groups {
			err := service.SendAnnouncement(group.ID, message)
			if err != nil {
				slog.Error("Failed to send announcement", "error", err, "group_id", group.ID)
			} else {
				slog.Info("Announcement sent", "group_id", group.ID)
				successCount++
			}
		}

		fmt.Printf("Announcement sent to %d/%d groups\n", successCount, len(groups))
	},
}

func init() {
	announceCmd.Flags().StringVarP(&announceFile, "file", "f", "", "Read message content from file")
	announceCmd.Flags().Int64VarP(&announceGroup, "group", "g", 0, "Send to specific group ID (default: all groups)")
	announceCmd.Flags().BoolVarP(&announcePreview, "preview", "p", false, "Preview message without sending")
}

