package ui

import (
	"github.com/eduqr/ktype/internal/typer"
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

func SetupStyles() tcell.Style {
	challengeStyle := tcell.StyleDefault.Foreground(color.Reset).Background(color.Reset)

	return challengeStyle
}

func DrawChallenge(screen tcell.Screen, cells []typer.Cell, style tcell.Style) {
	screen.Clear()

	for _, cell := range cells {
		screen.PutStrStyled(cell.X, cell.Y, cell.Expected, style)
	}
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

	screen.Show()
}
