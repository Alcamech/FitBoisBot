package bot

import (
	"log"
	"time"

	"github.com/Alcamech/FitBoisBot/internal/database/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func onFastGG(bot *tgbotapi.BotAPI, chatID int64) {
	group, err := groupRepo.GetGroupByID(chatID)
	if err != nil {
		log.Printf("Failed to fetch group %d: %v", chatID, err)
		sendText(bot, chatID, "Failed to fetch group timezone. Please set the timezone using /settimezone.")
		return
	}

	timezone := group.Timezone
	if timezone == "" {
		timezone = "America/New_York"
	}

	year, err := GetCurrentYear(timezone)
	if err != nil {
		log.Printf("Failed to get current year for timezone %s: %v", timezone, err)
		sendText(bot, chatID, "Invalid timezone set for this group. Please update the timezone using /settimezone.")
		return
	}

	leaderboard, err := ggRepo.GetGGLeaderboard(chatID, year)
	if err != nil {
		log.Printf("Failed to fetch GG leaderboard for chat %d: %v", chatID, err)
		return
	}

	message := formatFastGGLeaderboard(leaderboard)
	sendText(bot, chatID, message)
}

func onHelp(bot *tgbotapi.BotAPI, chatID int64) {
	message := "*🏋️ Welcome to the FitBois Bot\\!*\n" +
		"Here are the available commands:\n\n" +
		"*📜 Commands:*\n" +
		"  • `/help` \\- Show this help message\\.\n" +
		"  • `/fastgg` \\- Display the Fast GG leaderboard your group\\.\n" +
		"  • `/tokens` \\- Display Fitboi Token balances\\.\n" +
		"  • `/timezone` \\- Set the timezone for your group \\(e\\.g\\., `/timezone America/New_York`\\)\\.\n\n" +
		"*🗓️ Activity Format:*\n" +
		"Post activities in the format:\n" +
		"`activity\\-MM\\-DD\\-YYYY` or `activity\\-MM\\-DD\\-YY`\\.\n\n" +
		"*🚑 Support:*\n" +
		"For any issues, reach out to our support team\\.\n\n" +
		"💪 *Stay fit and keep posting your progress\\!*"

	sendMarkdownText(bot, chatID, message)
}

func onSetTimezone(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	args := msg.CommandArguments()
	if args == "" {
		sendText(bot, msg.Chat.ID, "Please specify a timezone. Example: /timezone America/New_York")
		return
	}

	_, err := time.LoadLocation(args) // Validate timezone
	if err != nil {
		sendText(bot, msg.Chat.ID, "Invalid timezone. Please use a valid IANA timezone. Example: /timezone America/New_York")
		return
	}

	err = groupRepo.SetGroupTimezone(msg.Chat.ID, args)
	if err != nil {
		log.Printf("Failed to set timezone for chat %d: %v", msg.Chat.ID, err)
		sendText(bot, msg.Chat.ID, "Failed to set timezone. Please try again.")
		return
	}

	sendText(bot, msg.Chat.ID, "Timezone updated successfully to "+args)
}

func onTokens(bot *tgbotapi.BotAPI, chatID int64) {
	group, err := groupRepo.GetGroupByID(chatID)
	if err != nil {
		log.Printf("Failed to fetch group %d: %v", chatID, err)
		sendText(bot, chatID, "Failed to fetch group timezone. Please set the timezone using /settimezone.")
		return
	}

	timezone := group.Timezone
	if timezone == "" {
		timezone = "America/New_York"
	}

	year, err := GetCurrentYear(timezone)
	if err != nil {
		log.Printf("Failed to get current year for timezone %s: %v", timezone, err)
		sendText(bot, chatID, "Invalid timezone set for this group. Please update the timezone using /settimezone.")
		return
	}

	leaderboard, err := tokenRepo.GetYearlyLeaderboard(chatID, year)
	if err != nil {
		log.Printf("Failed to fetch leaderboard for group %d: %v", chatID, err)
		sendText(bot, chatID, "Failed to fetch leaderboard. Please try again later.")
		return
	}

	message := formatTokenLeaderboard(leaderboard, year)
	sendText(bot, chatID, message)
}

func onActivityPost(msg *tgbotapi.Message, user *models.User) error {
	activity, month, day, year, err := parseActivityMessage(msg)
	if err != nil {
		sendReply(bot, msg.Chat.ID, "Wrong activity format. use activity-MM-DD-YYYY or activity-MM-DD-YY", msg.MessageID)
		return err
	}

	return activityRepo.CreateActivityRecord(user.ID, msg.Chat.ID, activity, month, day, year)
}
