package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type ChallengeConfig struct {
	Words  int    `mapstructure:"words"`
	Source string `mapstructure:"source"`
	Quote  bool   `mapstructure:"quote"`
	StopOn string `mapstructure:"stop_on"`
}

type LogsConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

type ThemeConfig struct {
	Normal    string `mapstructure:"normal"`
	Correct   string `mapstructure:"correct"`
	Incorrect string `mapstructure:"incorrect"`
	Missed    string `mapstructure:"missed"`
	Title     string `mapstructure:"title"`
	Accuracy  string `mapstructure:"accuracy"`
	WPM       string `mapstructure:"wpm"`
	RawWPM    string `mapstructure:"raw_wpm"`
	Time      string `mapstructure:"time"`
	Graph     string `mapstructure:"graph"`
	Actions   string `mapstructure:"actions"`
	Error     string `mapstructure:"error"`
}

type Config struct {
	Challenge ChallengeConfig `mapstructure:"challenge"`
	Logs      LogsConfig      `mapstructure:"logs"`
	Theme     ThemeConfig     `mapstructure:"theme"`
}

var Cfg Config

func Dir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine config directory: %w", err)
	}

	return filepath.Join(configDir, "ktype"), nil
}

func InitConfig() error {
	viper.SetDefault("challenge.words", 25)
	viper.SetDefault("challenge.source", "english")
	viper.SetDefault("challenge.quote", false)
	viper.SetDefault("challenge.stop_on", "off")
	viper.SetDefault("logs.enabled", false)
	viper.SetDefault("theme.normal", "#a9b1d6")
	viper.SetDefault("theme.correct", "#b9f27c")
	viper.SetDefault("theme.incorrect", "#f7768e")
	viper.SetDefault("theme.missed", "#853F4C")
	viper.SetDefault("theme.title", "#cfc9c2")
	viper.SetDefault("theme.accuracy", "#e0af68")
	viper.SetDefault("theme.wpm", "#7aa2f7")
	viper.SetDefault("theme.raw_wpm", "#7dcfff")
	viper.SetDefault("theme.time", "#bb9af7")
	viper.SetDefault("theme.graph", "#7aa2f7")
	viper.SetDefault("theme.actions", "#565f89")
	viper.SetDefault("theme.error", "#f7768e")

	ktypeDir, err := Dir()
	if err != nil {
		return err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ktypeDir)

	if err := os.MkdirAll(filepath.Join(ktypeDir, "words"), 0755); err != nil {
		return fmt.Errorf("cannot create words directory: %w", err)
	}

	if err := os.MkdirAll(filepath.Join(ktypeDir, "quotes"), 0755); err != nil {
		return fmt.Errorf("cannot create quotes directory: %w", err)
	}

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

	if err := viper.Unmarshal(&Cfg); err != nil {
		return err
	}

	if err := Cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	return nil
}

func (c Config) Validate() error {
	if err := validateStopOn(c.Challenge.StopOn); err != nil {
		return err
	}
	return nil
}

func validateStopOn(v string) error {
	switch v {
	case "off", "word", "letter":
		return nil
	default:
		return fmt.Errorf("challenge.stop_on: must be one of off, word, letter (got %q)", v)
	}
}
