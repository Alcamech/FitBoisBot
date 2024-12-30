package bot

import (
	"log"

	"github.com/Alcamech/FitBoisBot/internal/database/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func onFastGG(bot *tgbotapi.BotAPI, chatID int64) {
	leaderboard, err := ggRepo.GetGGLeaderboard(chatID)
	if err != nil {
		log.Printf("Failed to fetch GG leaderboard for chat %d: %v", chatID, err)
		return
	}

	message := formatFastGGLeaderboard(leaderboard)
	sendText(bot, chatID, message)
}

func onHelp(bot *tgbotapi.BotAPI, chatID int64) {
	message := "Nothing to see here... yet"
	sendText(bot, chatID, message)
}

func onActivityPost(msg *tgbotapi.Message, user *models.User) error {
	activity, month, day, year, err := parseActivityMessage(msg)
	if err != nil {
		sendReply(bot, msg.Chat.ID, "Wrong activity format. use activity-MM-DD-YYYY", msg.MessageID)
		return err
	}

	return activityRepo.CreateActivityRecord(user.ID, msg.Chat.ID, activity, month, day, year)
}
