package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppName string `mapstructure:"APP_NAME"`
	Port    int    `mapstructure:"PORT"`
	DB      struct {
		Host     string `mapstructure:"DB_HOST"`
		Port     string `mapstructure:"DB_PORT"`
		User     string `mapstructure:"DB_USER"`
		Password string `mapstructure:"DB_PASSWORD"`
		Name     string `mapstructure:"DB_NAME"`
		SSLMode  string `mapstructure:"DB_SSLMODE"`
	} `mapstructure:",squash"`
	CORS struct {
		AllowedOrigins string `mapstructure:"CORS_ALLOWED_ORIGINS"`
		AllowedMethods string `mapstructure:"CORS_ALLOWED_METHODS"`
		AllowedHeaders string `mapstructure:"CORS_ALLOWED_HEADERS"`
		ExposedHeaders string `mapstructure:"CORS_EXPOSED_HEADERS"`
	} `mapstructure:",squash"`
	ENCRYPTION_SECRET string
	AWS_REGION        string `mapstructure:"AWS_REGION"`
	SES_SENDER        string `mapstructure:"AWS_SES_SENDER"`
	// EMAIL_BACKEND can be "ses" or "local". When set to "local" emails are written to a file for testing.
	EMAIL_BACKEND string `mapstructure:"EMAIL_BACKEND"`
	// EMAIL_SAVE_PATH is used by the local backend to store generated emails (optional).
	EMAIL_SAVE_PATH string `mapstructure:"EMAIL_SAVE_PATH"`
	// QR HUBS
	HUBS string `mapstructure:"HUBS"`
	// SENDMAIL_PATH defaults to /usr/sbin/sendmail if empty
	SENDMAIL_PATH string `mapstructure:"SENDMAIL_PATH"`
}

var AppConfig Config

func LoadConfig(config_file string) {
	// Set defaults
	viper.SetDefault("APP_NAME", "p2p management service")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DB_SSLMODE", "disable")

	// Load .env file
	viper.SetConfigFile(config_file)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Error unmarshaling configuration: %v", err)
	}

	log.Println("Configuration loaded successfully!")
}
