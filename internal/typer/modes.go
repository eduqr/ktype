package typer

import (
	"slices"

	"github.com/eduqr/ktype/internal/ui"
)

func (gs *GameState) handleStopOnLetter(input string) Action {
	if input == " " && gs.currentCell().IsOverflow {
		gs.Cells = slices.Delete(gs.Cells, gs.CursorIndex, gs.CursorIndex+1)
	}

	cell := gs.currentCell()

	if input != cell.Expected {
		if cell.Expected == " " {
			gs.setOverflowCell(input)
		} else {
			cell.State = ui.StatusIncorrect
			gs.Errors++
		}

		return ActionNone
	}

	return gs.advance(true)
}

func (gs *GameState) handleStopOnOff(input string) Action {
	cell := gs.currentCell()

	if input == " " {
		if cell.Expected == " " {
			return gs.advance(true)
		}

		// the word is incomplete: skip what's left
		gs.skipRemainingCells()
		if gs.CursorIndex >= len(gs.Cells) {
			return gs.finish()
		}

		return ActionNone
	}

	// if the input wasn't space then it's overflow
	if cell.Expected == " " {
		if overflowCountForWord(gs.Cells, gs.CursorIndex) >= gs.Rules.MaxWordOverflow {
			return ActionNone
		}

		gs.appendOverflow(input)
		gs.CursorIndex++
		return ActionNone
	}

	return gs.advance(input == cell.Expected)
}

func (gs *GameState) handleStopOnWord(input string) Action {
	cell := gs.currentCell()

	if cell.Expected == " " && (input != " " || wordHasErrors(gs.Cells, gs.CursorIndex)) {
		if overflowCountForWord(gs.Cells, gs.CursorIndex) >= gs.Rules.MaxWordOverflow {
			return ActionNone
		}

		gs.appendOverflow(input)
		gs.CursorIndex++
		return ActionNone
	}

	return gs.advance(input == cell.Expected)
}

func (gs *GameState) skipRemainingCells() {
	word := gs.currentCell().WordIndex

	for gs.CursorIndex < len(gs.Cells) && gs.currentCell().WordIndex == word {
		cell := gs.currentCell()
		if cell.Expected == " " {
			cell.State = ui.StatusCorrect
		} else {
			cell.State = ui.StatusMissed
		}
		gs.CursorIndex++
	}

	gs.Errors++
}
