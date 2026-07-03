package main

import (
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	PORT      string `mapstructure:"PORT"`
	DB_URL    string `mapstructure:"DATABASE_URL"`
	JWTSecret string `mapstructure:"SECRET_KEY"`
}

func LoadConfig(logger *slog.Logger) (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error to load .env", "error", err)
		return nil, fmt.Errorf("Error to load .env %v", err)
	}

	viper.SetDefault("PORT", "8080")
	viper.BindEnv("DATABASE_URL")
	viper.BindEnv("SECRET_KEY")
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("Error to get variables")
	}

	if cfg.DB_URL == "" || cfg.JWTSecret == "" {
		logger.Error("Empty variables")
		return nil, fmt.Errorf("Empty variables")
	}

	return &cfg, nil
}
