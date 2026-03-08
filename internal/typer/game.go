package typer

import (
	"unicode/utf8"

	"github.com/eduqr/ktype/internal/ui"
	"github.com/gdamore/tcell/v3"
)

type GameState struct {
	Cells       []ui.Cell
	Styles      ui.Styles
	CursorIndex int
	TotalChars  int
	Errors      int
	IsDone      bool
}

func BuildCells(challenge []string, startX, startY int, screenWidth, screenHeight int) []ui.Cell {
	cells := make([]ui.Cell, 0)

	for i, word := range challenge {
		runeCount := utf8.RuneCountInString(word)
		if startX+runeCount >= screenWidth {
			startY++
			startX = 0
		}
		charPos := 0
		for _, char := range word {
			cell := ui.Cell{
				X:         startX + charPos,
				Y:         startY,
				Expected:  string(char),
				State:     0,
				WordIndex: i,
			}
			cells = append(cells, cell)
			charPos++
		}
		if i != len(challenge)-1 {
			cell := ui.Cell{
				X:         startX + runeCount,
				Y:         startY,
				Expected:  " ",
				State:     0,
				WordIndex: i,
			}
			cells = append(cells, cell)
		}
		startX += runeCount + 1
	}

	return cells
}

func InitGameState(cells []ui.Cell, styles ui.Styles) GameState {
	return GameState{
		Cells:       cells,
		Styles:      styles,
		CursorIndex: 0,
		TotalChars:  len(cells),
		Errors:      0,
		IsDone:      false,
	}
}

func StartGame(screen tcell.Screen, gameState GameState) {
	drawCursor(screen, gameState)

	for {
		event := <-screen.EventQ()

		switch ev := event.(type) {
		case *tcell.EventResize:
			screen.Sync()
			ui.DrawChallenge(screen, gameState.Cells, gameState.Styles)
			drawCursor(screen, gameState)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			}

			if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
				if gameState.CursorIndex <= 0 || gameState.IsDone {
					continue
				}

				if gameState.Cells[gameState.CursorIndex].WordIndex != gameState.Cells[gameState.CursorIndex-1].WordIndex {
					continue
				}

				ui.DrawCell(screen, gameState.Cells[gameState.CursorIndex], gameState.Styles)
				gameState.CursorIndex--
				if gameState.Cells[gameState.CursorIndex].State == 2 {
					gameState.Errors--
				}
				gameState.Cells[gameState.CursorIndex].State = 0
				ui.DrawCursorCell(screen, gameState.Cells[gameState.CursorIndex], gameState.Styles)

			}
			if ev.Key() == tcell.KeyRune {
				if gameState.IsDone {
					return
				}

				input := ev.Str()
				currentIndex := gameState.CursorIndex
				expected := gameState.Cells[currentIndex].Expected

				if input == expected {
					gameState.Cells[currentIndex].State = 1
				} else if expected == " " {
					gameState.Cells[currentIndex].State = 1
				} else {
					gameState.Cells[currentIndex].State = 2
					gameState.Errors++
				}

				ui.DrawCell(screen, gameState.Cells[currentIndex], gameState.Styles)
				gameState.CursorIndex++

				if gameState.CursorIndex >= gameState.TotalChars {
					gameState.IsDone = true
					screen.Show()
					return
				}

				ui.DrawCursorCell(screen, gameState.Cells[gameState.CursorIndex], gameState.Styles)
			}
		}

		drawCursor(screen, gameState)
		screen.Show()
	}
}

func drawCursor(screen tcell.Screen, gameState GameState) {
	if gameState.CursorIndex >= gameState.TotalChars {
		return
	}

	cell := gameState.Cells[gameState.CursorIndex]
	ui.DrawCursorCell(screen, cell, gameState.Styles)
	screen.ShowCursor(cell.X, cell.Y)
}
