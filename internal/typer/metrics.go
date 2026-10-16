package typer

import "github.com/eduqr/ktype/internal/ui"

func (gs *GameState) countCorrectCells() int {
	count := 0
	for _, cell := range gs.Cells {
		if cell.State == ui.StatusCorrect && !cell.IsTerminator {
			count++
		}
	}
	return count
}
