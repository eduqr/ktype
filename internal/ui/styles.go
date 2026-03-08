package ui

import (
	"fmt"
	"strings"

	"github.com/eduqr/ktype/internal/config"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Styles struct {
	Normal           tcell.Style
	Correct          tcell.Style
	Incorrect        tcell.Style
	cursorBackground color.Color
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

	cursorColor, err := resolveThemeColor("theme.cursor", theme.Cursor)
	if err != nil {
		return Styles{}, err
	}

	return Styles{
		Normal:           tcell.StyleDefault.Foreground(normalColor).Background(color.Reset),
		Correct:          tcell.StyleDefault.Foreground(correctColor).Background(color.Reset),
		Incorrect:        tcell.StyleDefault.Foreground(incorrectColor).Background(color.Reset),
		cursorBackground: cursorColor,
	}, nil
}

func (styles Styles) CellStyle(state int) tcell.Style {
	switch state {
	case 1:
		return styles.Correct
	case 2:
		return styles.Incorrect
	default:
		return styles.Normal
	}
}

func (styles Styles) CursorStyle(state int) tcell.Style {
	return styles.CellStyle(state).Background(styles.cursorBackground)
}

func resolveThemeColor(fieldName, value string) (color.Color, error) {
	parsedColor := color.GetColor(strings.TrimSpace(value))
	if !parsedColor.Valid() {
		return color.Default, fmt.Errorf("invalid color for %s: %q", fieldName, value)
	}

	return parsedColor, nil
}
