package main

import (
	"github.com/Alcamech/FitBoisBot/config"
	"github.com/Alcamech/FitBoisBot/internal/bot"
	"github.com/Alcamech/FitBoisBot/internal/database"
)

func main() {
	config.InitConfig()
	database.InitDB()
	bot.StartBot()
}
