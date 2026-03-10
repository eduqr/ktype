package ui

import "unicode/utf8"

type CellStatus int

const (
	StatusPending CellStatus = iota
	StatusCorrect
	StatusIncorrect
)

type Cell struct {
	X, Y      int
	Expected  string
	State     CellStatus
	WordIndex int
}

func BuildCells(challenge []string, startX, startY int, screenWidth int) []Cell {
	cells := make([]Cell, 0)

	for i, word := range challenge {
		runeCount := utf8.RuneCountInString(word)
		if startX > 0 && startX+runeCount >= screenWidth {
			startY++
			startX = 0
		}
		charPos := 0
		for _, char := range word {
			cell := Cell{
				X:         startX + charPos,
				Y:         startY,
				Expected:  string(char),
				State:     StatusPending,
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
				State:     StatusPending,
				WordIndex: i,
			}
			cells = append(cells, cell)
		}
		startX += runeCount + 1
	}

	return cells
}
