package ui

import (
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
	challengeStyle := tcell.StyleDefault.Foreground(color.YellowGreen).Background(color.Reset)

	return challengeStyle
}

func DrawChallenge(screen tcell.Screen, challenge []string, style tcell.Style) {
	screen.Clear()

	width, _ := screen.Size()

	x, y := 0, 1
	for _, word := range challenge {

		if x+len(word) >= width {
			x = 0
			y++
		}

		screen.PutStrStyled(x, y, word, style)
		x += len(word)

		screen.Put(x, y, " ", style)
		x++
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
