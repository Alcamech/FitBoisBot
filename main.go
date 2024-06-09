package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

func main() {
	initConfig()

	fmt.Println("Hello, World")

	fmt.Println("Database URL:", viper.GetString("database.url"))
}
