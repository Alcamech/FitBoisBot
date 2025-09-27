package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Username string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	Name     string `mapstructure:"DB_NAME"`
	Token    string `mapstructure:"TOKEN"`
	Debug    bool   `mapstructure:"DEBUG"`
}

var AppConfig Config

func InitConfig() {
	viper.AutomaticEnv()

	// Read .env file if it exists
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env file found, using environment variables and defaults: %v", err)
	} else {
		log.Printf("Loaded .env file: %s", viper.ConfigFileUsed())
	}

	viper.SetDefault("DEBUG", false)

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	// Validate required configuration
	if AppConfig.Token == "" {
		log.Fatalf("Telegram token is required (set TOKEN environment variable)")
	}

	log.Printf("Configuration loaded successfully")
}
