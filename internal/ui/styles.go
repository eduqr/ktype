package ui

import (
	"fmt"
	"strings"

	"github.com/eduqr/ktype/internal/config"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Styles struct {
	Normal    tcell.Style
	Correct   tcell.Style
	Incorrect tcell.Style
	Missed    tcell.Style
	Title     tcell.Style
	Accuracy  tcell.Style
	WPM       tcell.Style
	RawWPM    tcell.Style
	Time      tcell.Style
	Graph     tcell.Style
	Actions   tcell.Style
	Error     tcell.Style
}

func NewStyles(theme config.ThemeConfig) (Styles, error) {
	normalColor, err := resolveThemeColor("theme.normal", theme.Normal)
	if err != nil {
		return Styles{}, err
	}

	correctColor, err := resolveThemeColor("theme.correct", theme.Correct)
	if err != nil {
		return Styles{}, err
	}

	incorrectColor, err := resolveThemeColor("theme.incorrect", theme.Incorrect)
	if err != nil {
		return Styles{}, err
	}

	missedColor, err := resolveThemeColor("theme.missed", theme.Missed)
	if err != nil {
		return Styles{}, err
	}

	titleColor, err := resolveThemeColor("theme.title", theme.Title)
	if err != nil {
		return Styles{}, err
	}

	accuracyColor, err := resolveThemeColor("theme.accuracy", theme.Accuracy)
	if err != nil {
		return Styles{}, err
	}

	wpmColor, err := resolveThemeColor("theme.wpm", theme.WPM)
	if err != nil {
		return Styles{}, err
	}

	rawWPMColor, err := resolveThemeColor("theme.raw_wpm", theme.RawWPM)
	if err != nil {
		return Styles{}, err
	}

	timeColor, err := resolveThemeColor("theme.time", theme.Time)
	if err != nil {
		return Styles{}, err
	}

	graphColor, err := resolveThemeColor("theme.graph", theme.Graph)
	if err != nil {
		return Styles{}, err
	}

	actionsColor, err := resolveThemeColor("theme.actions", theme.Actions)
	if err != nil {
		return Styles{}, err
	}

	errorColor, err := resolveThemeColor("theme.error", theme.Error)
	if err != nil {
		return Styles{}, err
	}

	return Styles{
		Normal:    tcell.StyleDefault.Foreground(normalColor).Background(color.Reset),
		Correct:   tcell.StyleDefault.Foreground(correctColor).Background(color.Reset),
		Incorrect: tcell.StyleDefault.Foreground(incorrectColor).Background(color.Reset),
		Missed:    tcell.StyleDefault.Foreground(missedColor).Background(color.Reset),
		Title:     tcell.StyleDefault.Foreground(titleColor).Background(color.Reset),
		Accuracy:  tcell.StyleDefault.Foreground(accuracyColor).Background(color.Reset),
		WPM:       tcell.StyleDefault.Foreground(wpmColor).Background(color.Reset),
		RawWPM:    tcell.StyleDefault.Foreground(rawWPMColor).Background(color.Reset),
		Time:      tcell.StyleDefault.Foreground(timeColor).Background(color.Reset),
		Graph:     tcell.StyleDefault.Foreground(graphColor).Background(color.Reset),
		Actions:   tcell.StyleDefault.Foreground(actionsColor).Background(color.Reset),
		Error:     tcell.StyleDefault.Foreground(errorColor).Background(color.Reset),
	}, nil
}

func (styles Styles) cellStyle(state CellStatus) tcell.Style {
	switch state {
	case StatusCorrect:
		return styles.Correct
	case StatusIncorrect:
		return styles.Incorrect
	case StatusMissed:
		return styles.Missed
	default:
		return styles.Normal
	}
}

func resolveThemeColor(fieldName, value string) (color.Color, error) {
	parsedColor := color.GetColor(strings.TrimSpace(value))
	if !parsedColor.Valid() {
		return color.Default, fmt.Errorf("invalid color for %s: %q", fieldName, value)
	}

	return parsedColor, nil
}
