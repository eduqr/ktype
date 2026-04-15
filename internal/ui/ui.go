package ui

import (
	"fmt"

	"github.com/eduqr/ktype/internal/stats"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func InitScreen() (tcell.Screen, error) {
	defStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)

	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	screen.SetStyle(defStyle)
	screen.DisablePaste()
	screen.SetTitle("ktype")
	screen.Clear()

	return screen, nil
}

func DrawChallenge(screen tcell.Screen, cells []Cell, styles Styles) {
	screen.Clear()

	for _, cell := range cells {
		DrawCell(screen, cell, styles)
	}
}

func DrawCell(screen tcell.Screen, cell Cell, styles Styles) {
	screen.Put(cell.X, cell.Y, cell.Expected, styles.CellStyle(cell.State))
}

func DrawCursorCell(screen tcell.Screen, cell Cell, styles Styles) {
	screen.Put(cell.X, cell.Y, cell.Expected, styles.CursorStyle(cell.State))
}

func ShowFinalScreen(screen tcell.Screen, styles Styles, stats stats.Stats) {
	screen.Clear()
	screen.HideCursor()

	drawTitle(screen, styles.Normal)
	drawResults(screen, stats, styles)
	showOptions(screen, styles)

	screen.Show()
}

func drawResults(screen tcell.Screen, stats stats.Stats, styles Styles) {
	result := fmt.Sprintf("Acc: %.2f, WPM: %.2f, Raw WPM: %.2f", stats.Accuracy, stats.WPM, stats.RawWPM)

	screen.PutStrStyled(1, 8, result, styles.Normal)
}

func showOptions(screen tcell.Screen, styles Styles) {
	screen.PutStrStyled(1, 11, "[r] repeat", styles.Normal)
	screen.PutStrStyled(1, 12, "[tab] new challenge", styles.Normal)
	screen.PutStrStyled(1, 13, "[esc] quit", styles.Normal)
}

func drawTitle(screen tcell.Screen, style tcell.Style) {
	title := []string{
		" _    _                    ",
		"| | _| |_ _   _ _ __   ___ ",
		"| |/ / __| | | | '_ \\ / _ \\",
		"|   <| |_| |_| | |_) |  __/",
		"|_|\\_\\\\__|\\__, | .__/ \\___|",
		"          |___/|_|         ",
	}

	x, y := 1, 1

	for _, line := range title {
		screen.PutStrStyled(x, y, line, style)
		y++
	}
}
