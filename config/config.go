package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Name     string `mapstructure:"name"`
	}

	Telegram struct {
		DevToken  string `mapstructure:"dev-token"`
		QaToken   string `mapstructure:"qa-token"`
		ProdToken string `mapstructure:"prod-token"`
		Debug     bool   `mapstructure:"debug-mode"`
	}
}

var AppConfig Config

func InitConfig() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("fitbois")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found or not used: %v", err)
	}

	viper.SetDefault("telegram.debug-mode", false)

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	log.Printf("Loaded Config: %+v", AppConfig)
}
