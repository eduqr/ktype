package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type StopOn string

const (
	StopOnOff    StopOn = "off"
	StopOnWord   StopOn = "word"
	StopOnLetter StopOn = "letter"
)

type ChallengeConfig struct {
	Words  int    `mapstructure:"words"`
	Source string `mapstructure:"source"`
	Quote  bool   `mapstructure:"quote"`
	StopOn StopOn `mapstructure:"stop_on"`
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
	setDefaults(viper.GetViper())

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
		defaultConfig := viper.New()
		setDefaults(defaultConfig)
		if err := defaultConfig.SafeWriteConfigAs(filepath.Join(ktypeDir, "config.yaml")); err != nil {
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

func setDefaults(v *viper.Viper) {
	v.SetDefault("challenge.words", 25)
	v.SetDefault("challenge.source", "english")
	v.SetDefault("challenge.quote", false)
	v.SetDefault("challenge.stop_on", StopOnOff)
	v.SetDefault("theme.normal", "#a9b1d6")
	v.SetDefault("theme.correct", "#b9f27c")
	v.SetDefault("theme.incorrect", "#f7768e")
	v.SetDefault("theme.missed", "#853F4C")
	v.SetDefault("theme.title", "#cfc9c2")
	v.SetDefault("theme.accuracy", "#e0af68")
	v.SetDefault("theme.wpm", "#7aa2f7")
	v.SetDefault("theme.raw_wpm", "#7dcfff")
	v.SetDefault("theme.time", "#bb9af7")
	v.SetDefault("theme.graph", "#7aa2f7")
	v.SetDefault("theme.actions", "#565f89")
	v.SetDefault("theme.error", "#f7768e")
}

func (c Config) Validate() error {
	if err := validateStopOn(c.Challenge.StopOn); err != nil {
		return err
	}
	return nil
}

func validateStopOn(v StopOn) error {
	switch v {
	case StopOnOff, StopOnWord, StopOnLetter:
		return nil
	default:
		return fmt.Errorf("challenge.stop_on: must be one of off, word, letter (got %q)", v)
	}
}
