package ui

type CellStatus int

const (
	StatusPending CellStatus = iota
	StatusCorrect
	StatusIncorrect
	StatusMissed
)

type Cell struct {
	X, Y         int
	Expected     string
	OverflowChar string
	State        CellStatus
	WordIndex    int
	IsOverflow   bool
	IsTerminator bool
}

func BuildCells(challenge []string, layout Layout) []Cell {
	cells := make([]Cell, 0)

	for i, word := range challenge {
		for _, char := range word {
			cells = append(cells, Cell{
				Expected:  string(char),
				State:     StatusPending,
				WordIndex: i,
			})
		}
		// last one = terminator
		cells = append(cells, Cell{
			Expected:     " ",
			State:        StatusPending,
			WordIndex:    i,
			IsTerminator: i == len(challenge)-1,
		})
	}

	RecalcChallengePositions(cells, layout)
	return cells
}

func RecalcChallengePositions(cells []Cell, layout Layout) {
	RecalcPositions(cells, layout.Left, 0, layout.ContentWidth)
	lineCount := cellLineCount(cells)
	startY := challengeStartY(layout, lineCount)

	RecalcPositions(cells, layout.Left, startY, layout.ContentWidth)
}

func cellLineCount(cells []Cell) int {
	if len(cells) == 0 {
		return 0
	}

	maxY := cells[0].Y
	for _, cell := range cells[1:] {
		if cell.Y > maxY {
			maxY = cell.Y
		}
	}

	return maxY + 1
}

func challengeStartY(layout Layout, lineCount int) int {
	if lineCount <= 0 {
		return layout.Top
	}

	top := max((layout.ScreenHeight-lineCount)/4, layout.Top)
	if top+lineCount > layout.ScreenHeight {
		top = max(layout.ScreenHeight-lineCount, 0)
	}

	return top
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

func RecalcPositions(cells []Cell, startX, startY, lineWidth int) {
	if lineWidth < 1 {
		lineWidth = 1
	}

	x, y := startX, startY
	rightX := startX + lineWidth

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
			if isWordStart && x > startX {
				wordLen := wordGroupLen(cells, i)
				if x+wordLen > rightX {
					x = startX
					y++
				}
			}
		}

		cell.X = x
		cell.Y = y
		x++

		if x >= rightX {
			x = startX
			y++
		}
	}
}
