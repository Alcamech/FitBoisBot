package main

import (
	"log"

	"github.com/Alcamech/FitBoisBot/config"
	"github.com/Alcamech/FitBoisBot/internal/bot"
	"github.com/Alcamech/FitBoisBot/internal/database"
)

func main() {
	config.InitConfig()
	database.InitDB()

	go func() {
		log.Println("Starting monthly token awards scheduler...")
		bot.ScheduleMonthlyTokenAwards()
	}()

	bot.BotLoop()
}
