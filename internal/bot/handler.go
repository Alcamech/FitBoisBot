package bot

import (
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func onFastGG(bot *tgbotapi.BotAPI, chatID int64) {
	message := "Fastest GG: (Sample Data)"
	sendText(bot, chatID, message)
}

func onHelp(bot *tgbotapi.BotAPI, chatID int64) {
	message := "Nothing to see here... yet"
	sendText(bot, chatID, message)
}

func onActivityPost(msg *tgbotapi.Message, user *models.User) error {
	activity, month, day, year, err := parseActivityMessage(msg)
	if err != nil {
		return err
	}

	return activityRepo.CreateActivityRecord(user.ID, activity, month, day, year)
}
