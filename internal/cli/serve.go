package cli

import (
	"log/slog"

	"github.com/Alcamech/FitBoisBot/config"
	"github.com/Alcamech/FitBoisBot/internal/bot"
	"github.com/Alcamech/FitBoisBot/internal/database"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the bot server",
	Long:  `Start the FitBoisBot Telegram bot server with message processing and monthly token scheduler.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.InitConfig()
		database.InitDB()

		schedulerService, err := bot.NewBotService()
		if err != nil {
			slog.Error("Failed to initialize scheduler service", "error", err)
			panic(err)
		}

		go func() {
			slog.Info("Starting monthly token awards scheduler")
			schedulerService.ScheduleMonthlyTokenAwards()
		}()

		slog.Info("Starting FitBoisBot server")
		bot.BotLoop()
	},
}

