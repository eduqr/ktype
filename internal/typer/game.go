package typer

import (
	"unicode/utf8"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Cell struct {
	X, Y      int
	Expected  string
	State     int
	WordIndex int
}

type GameState struct {
	Cells       []Cell
	CursorIndex int
	TotalChars  int
	Errors      int
	IsDone      bool
}

func BuildCells(challenge []string, startX, startY int, screenWidth, screenHeight int) []Cell {
	cells := make([]Cell, 0)

	for i, word := range challenge {
		runeCount := utf8.RuneCountInString(word)
		if startX+runeCount >= screenWidth {
			startY++
			startX = 0
		}
		charPos := 0
		for _, char := range word {
			cell := Cell{
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
			cell := Cell{
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

func InitGameState(cells []Cell) GameState {
	return GameState{
		Cells:       cells,
		CursorIndex: 0,
		TotalChars:  len(cells),
		Errors:      0,
		IsDone:      false,
	}
}

func StartGame(screen tcell.Screen, gameState GameState) {
	for {
		screen.Show()
		event := <-screen.EventQ()

		switch ev := event.(type) {
		case *tcell.EventResize:
			screen.Sync()
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

				gameState.CursorIndex--
				cell := gameState.Cells[gameState.CursorIndex]
				screen.Put(cell.X, cell.Y, cell.Expected, tcell.StyleDefault)
				gameState.Cells[gameState.CursorIndex].State = 0
				if cell.State == 2 {
					gameState.Errors--
				}

			}
			if ev.Key() == tcell.KeyRune {
				if gameState.IsDone {
					return
				}

				input := ev.Str()
				cell := gameState.Cells[gameState.CursorIndex]
				expected := cell.Expected
				x := cell.X
				y := cell.Y

				if input == expected {
					if expected != " " {
						screen.Put(x, y, expected, tcell.StyleDefault.Foreground(color.Green))
					}
					gameState.Cells[gameState.CursorIndex].State = 1
				} else if expected == " " {
					screen.Put(x, y, "_", tcell.StyleDefault.Foreground(color.Red))
					gameState.Cells[gameState.CursorIndex].State = 1
				} else {
					screen.Put(x, y, expected, tcell.StyleDefault.Foreground(color.Red))
					gameState.Cells[gameState.CursorIndex].State = 2
					gameState.Errors++
				}
				gameState.CursorIndex++

				if gameState.CursorIndex >= gameState.TotalChars {
					gameState.IsDone = true
					return
				}
			}
		}

		screen.ShowCursor(gameState.Cells[gameState.CursorIndex].X, gameState.Cells[gameState.CursorIndex].Y)

	}
}
