package typer

import (
	"slices"

	"github.com/eduqr/ktype/internal/ui"
)

func (gs *GameState) canAppendOverflow() bool {
	return overflowCountForWord(gs.Cells, gs.CursorIndex) < gs.Rules.MaxWordOverflow
}

func (gs *GameState) appendOverflow(input string) {
	overflowCell := ui.Cell{
		Expected:     " ", // always space for overflow cells
		State:        ui.StatusIncorrect,
		WordIndex:    gs.Cells[gs.CursorIndex].WordIndex,
		OverflowChar: overflowChar(input),
		IsOverflow:   true,
	}
	gs.Cells = slices.Insert(gs.Cells, gs.CursorIndex, overflowCell)
	gs.Errors++
}

func (gs *GameState) setOverflowCell(input string) {
	if gs.Cells[gs.CursorIndex].IsOverflow {
		gs.Cells[gs.CursorIndex].OverflowChar = overflowChar(input)
		gs.Errors++
		return
	}

	gs.appendOverflow(input)
}

func (gs *GameState) removeCurrentOverflow() {
	gs.Cells = slices.Delete(gs.Cells, gs.CursorIndex, gs.CursorIndex+1)
}

func overflowChar(input string) string {
	if input == " " {
		return "_"
	}
	return input
}
