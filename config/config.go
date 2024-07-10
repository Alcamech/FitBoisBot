package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Telegram struct {
		DevToken  string `mapstructure:"dev-token"`
		QaToken   string `mapstructure:"qa-token"`
		ProdToken string `mapstructure:"prod-token"`
	}

	Database struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Name     string `mapstructure:"name"`
	}
}

var AppConfig Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}
