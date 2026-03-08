package ui

import (
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Cell struct {
	X, Y      int
	Expected  string
	State     int
	WordIndex int
}

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

func DrawChallenge(screen tcell.Screen, cells []Cell, style tcell.Style) {
	screen.Clear()

	for _, cell := range cells {
		screen.PutStrStyled(cell.X, cell.Y, cell.Expected, style)
	}
}

func DrawCell(screen tcell.Screen, cell Cell) {
	var cellStyle tcell.Style
	switch cell.State {
	case 0:
		cellStyle = tcell.StyleDefault.Foreground(color.Reset).Background(color.Reset)
	case 1:
		cellStyle = tcell.StyleDefault.Foreground(color.Green).Background(color.Reset)
	case 2:
		cellStyle = tcell.StyleDefault.Foreground(color.Red).Background(color.Reset)
	}

	screen.Put(cell.X, cell.Y, cell.Expected, cellStyle)
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
