package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type MailConfig struct {
	Port                string
	Host                string
	EmailSenderAddress  string
	EmailSenderPassword string
}

func New() (MailConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return MailConfig{}, fmt.Errorf("failed to load environment file: %w", err)
	}

	config := MailConfig{
		Port:                os.Getenv("SMTP_ADDR"),
		Host:                os.Getenv("FROM_EMAIL_SMTP"),
		EmailSenderAddress:  os.Getenv("FROM_EMAIL"),
		EmailSenderPassword: os.Getenv("FROM_EMAIL_PASSWORD"),
	}

	// Check if all values are loaded
	if config.Port == "" || config.Host == "" || config.EmailSenderAddress == "" || config.EmailSenderPassword == "" {
		return MailConfig{}, errors.New("incomplete mail configuration")
	}

	return config, nil
}
