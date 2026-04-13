package ui

import (
	"fmt"

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

func ShowResults(screen tcell.Screen, accuracy, wpm, rawWPM float64, styles Styles) {
	screen.Clear()
	screen.HideCursor()

	result := fmt.Sprintf("Acc: %.2f, WPM: %.2f, Raw WPM: %.2f", accuracy, wpm, rawWPM)

	drawTitle(screen, styles.Normal)
	screen.PutStrStyled(1, 8, result, styles.Normal)
	screen.Show()
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
