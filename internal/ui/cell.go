package ui

type CellStatus int

const (
	StatusPending CellStatus = iota
	StatusCorrect
	StatusIncorrect
)

type Cell struct {
	X, Y       int
	Expected   string
	State      CellStatus
	WordIndex  int
	IsOverflow bool
}

func BuildCells(challenge []string, startX, startY int, screenWidth int) []Cell {
	cells := make([]Cell, 0)

	for i, word := range challenge {
		for _, char := range word {
			cell := Cell{
				Expected:  string(char),
				State:     StatusPending,
				WordIndex: i,
			}
			cells = append(cells, cell)
		}
		if i != len(challenge)-1 {
			cell := Cell{
				Expected:  " ",
				State:     StatusPending,
				WordIndex: i,
			}
			cells = append(cells, cell)
		}
	}

	RecalcPositions(cells, startX, startY, screenWidth)
	return cells
}

// returns the number of cells belonging to the same word index (word + space + overflow)
func wordGroupLen(cells []Cell, i int) int {
	wordIndex := cells[i].WordIndex
	n := 0
	for j := i; j < len(cells) && cells[j].WordIndex == wordIndex; j++ {
		n++
	}
	return n
}

func RecalcPositions(cells []Cell, startX, startY, screenWidth int) {
	x, y := startX, startY

	for i := range cells {
		cell := &cells[i]

		if cell.Expected != " " && !cell.IsOverflow {
			isWordStart := false
			if i == 0 {
				isWordStart = true
			} else {
				prev := &cells[i-1]
				if prev.WordIndex != cell.WordIndex {
					isWordStart = true
				}
			}
			if isWordStart && x > 0 {
				wordLen := wordGroupLen(cells, i)
				if x+wordLen > screenWidth {
					x = 0
					y++
				}
			}
		}

		cell.X = x
		cell.Y = y
		x++

		if x >= screenWidth {
			x = 0
			y++
		}
	}
}
