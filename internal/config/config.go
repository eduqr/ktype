package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type ChallengeConfig struct {
	Words    int    `mapstructure:"words"`
	Language string `mapstructure:"language"`
	Quote    bool   `mapstructure:"quote"`
}

type LogsConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

type ThemeConfig struct {
	Normal    string `mapstructure:"normal"`
	Correct   string `mapstructure:"correct"`
	Incorrect string `mapstructure:"incorrect"`
	Cursor    string `mapstructure:"cursor"`
}

type Config struct {
	Challenge ChallengeConfig `mapstructure:"challenge"`
	Logs      LogsConfig      `mapstructure:"logs"`
	Theme     ThemeConfig     `mapstructure:"theme"`
}

var Cfg Config

func InitConfig() error {
	viper.SetDefault("challenge.words", 25)
	viper.SetDefault("challenge.language", "english")
	viper.SetDefault("challenge.quote", false)
	viper.SetDefault("logs.enabled", false)
	viper.SetDefault("theme.normal", "#ffffff")
	viper.SetDefault("theme.correct", "#00ff00")
	viper.SetDefault("theme.incorrect", "#ff0000")
	viper.SetDefault("theme.cursor", "#ffff00")

	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("cannot determine config directory: %w", err)
	}
	ktypeDir := filepath.Join(configDir, "ktype")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ktypeDir)

	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return fmt.Errorf("error reading config: %w", err)
		}
		if err := os.MkdirAll(ktypeDir, 0755); err != nil {
			return fmt.Errorf("cannot create config directory: %w", err)
		}
		if err := viper.SafeWriteConfigAs(filepath.Join(ktypeDir, "config.yaml")); err != nil {
			return fmt.Errorf("cannot write default config: %w", err)
		}
	}

	return viper.Unmarshal(&Cfg)
}
