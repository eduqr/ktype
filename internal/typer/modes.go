package typer

import (
	"github.com/eduqr/ktype/internal/ui"
)

func (gs *GameState) handleStopOnLetter(input string) Action {
	if input == " " && gs.currentCell().IsOverflow {
		gs.removeCurrentOverflow()
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
		if !gs.canAppendOverflow() {
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
		if !gs.canAppendOverflow() {
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

func (gs *GameState) restoreSkippedWord() bool {
	previousWord := gs.Cells[gs.CursorIndex-1].WordIndex
	start := wordStartIndex(gs.Cells, gs.CursorIndex-1)
	firstMissed := -1

	for i := start; i < gs.CursorIndex && gs.Cells[i].WordIndex == previousWord; i++ {
		if gs.Cells[i].State == ui.StatusMissed {
			if firstMissed == -1 {
				firstMissed = i
			}
			gs.Cells[i].State = ui.StatusPending
		}
	}

	if firstMissed == -1 {
		return false
	}

	gs.Cells[gs.CursorIndex-1].State = ui.StatusPending
	gs.CursorIndex = firstMissed
	return true
}
