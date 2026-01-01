package bot

import (
	"fmt"

	"github.com/Alcamech/FitBoisBot/internal/constants"
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// createChallengeListKeyboard creates an inline keyboard for challenge list pagination.
func createChallengeListKeyboard(challenges []models.Challenge, page, totalPages int) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Navigation row
	var navRow []tgbotapi.InlineKeyboardButton

	if page > 1 {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData(
			constants.ButtonPrevious,
			fmt.Sprintf("challenge_list_%d", page-1),
		))
	}

	if page < totalPages {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData(
			constants.ButtonNext,
			fmt.Sprintf("challenge_list_%d", page+1),
		))
	}

	if len(navRow) > 0 {
		rows = append(rows, navRow)
	}

	// View buttons for each challenge (up to 5 per row, max 2 rows)
	var viewRow []tgbotapi.InlineKeyboardButton
	for i, c := range challenges {
		viewRow = append(viewRow, tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d", i+1),
			fmt.Sprintf("challenge_view_%d", c.ID),
		))

		// Create a new row every 5 buttons or at the end
		if len(viewRow) == 5 || i == len(challenges)-1 {
			rows = append(rows, viewRow)
			viewRow = []tgbotapi.InlineKeyboardButton{}
		}
	}

	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}
