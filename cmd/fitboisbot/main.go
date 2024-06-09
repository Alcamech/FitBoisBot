package main

import (
	"github.com/Alcamech/FitBoisBot/config"
	"github.com/Alcamech/FitBoisBot/internal/database"
	"log"
)

func main() {
	config.InitConfig()
	database.InitDB()

	log.Println("FitBoisBot started")
}
