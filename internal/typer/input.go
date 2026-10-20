package typer

import (
	"slices"
	"time"

	"github.com/eduqr/ktype/internal/config"
	"github.com/eduqr/ktype/internal/ui"
	"github.com/gdamore/tcell/v3"
)

type Action int

const (
	ActionNone Action = iota
	ActionQuit
	ActionDone
	ActionAbortRound
	ActionRepeat
	ActionNewChallenge
)

func (gs *GameState) handleKey(ev *tcell.EventKey) Action {
	if ev.Key() == tcell.KeyCtrlC {
		return ActionQuit
	}

	if ev.Key() == tcell.KeyEscape {
		return ActionAbortRound
	}

	if !gs.layout.IsSupported() {
		return ActionNone
	}

	if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
		gs.handleBackspace()
		return ActionNone
	}

	if ev.Key() == tcell.KeyRune {
		return gs.handleRune(ev.Str())
	}

	return ActionNone
}

func (gs *GameState) handleBackspace() {
	if gs.CursorIndex <= 0 || gs.IsDone {
		return
	}

	stopOnMode := gs.Rules.StopOn

	// backspace right after a skip restores the skipped word
	if stopOnMode == config.StopOnOff &&
		isIndexAtWordStart(gs.Cells, gs.CursorIndex) &&
		gs.restoreSkippedWord() {
		return
	}

	// a wrong letter stays under the cursor so we have to clear the cell first
	if stopOnMode == config.StopOnLetter &&
		gs.Cells[gs.CursorIndex].State == ui.StatusIncorrect &&
		!gs.Cells[gs.CursorIndex].IsOverflow &&
		!isIndexAtWordStart(gs.Cells, gs.CursorIndex) {
		gs.Cells[gs.CursorIndex].State = ui.StatusPending
	}

	if stopOnMode == config.StopOnLetter && gs.Cells[gs.CursorIndex].IsOverflow {
		gs.Cells = slices.Delete(gs.Cells, gs.CursorIndex, gs.CursorIndex+1)
		gs.CursorIndex--
		gs.Cells[gs.CursorIndex].State = ui.StatusPending
		return
	}

	if gs.Cells[gs.CursorIndex-1].IsOverflow {
		gs.Cells = slices.Delete(gs.Cells, gs.CursorIndex-1, gs.CursorIndex)
		gs.CursorIndex--
		return
	}

	// block crossing into the prev word if it was completed correctly
	if gs.Cells[gs.CursorIndex].WordIndex != gs.Cells[gs.CursorIndex-1].WordIndex {
		if isPrevWordCorrect(gs.Cells, gs.CursorIndex) {
			return
		}
	}

	gs.CursorIndex--
	gs.Cells[gs.CursorIndex].State = ui.StatusPending
}

func (gs *GameState) handleRune(input string) Action {
	if gs.IsDone {
		return ActionDone
	}

	if gs.StartTime.IsZero() && gs.CursorIndex == 0 && input != " " {
		gs.StartTime = time.Now()
		gs.startWPMSampler()
	}

	// a space at the start of word does nothing
	if isIndexAtWordStart(gs.Cells, gs.CursorIndex) && input == " " {
		return ActionNone
	}

	switch gs.Rules.StopOn {
	case config.StopOnLetter:
		return gs.handleStopOnLetter(input)
	case config.StopOnOff:
		return gs.handleStopOnOff(input)
	case config.StopOnWord:
		return gs.handleStopOnWord(input)
	default:
		panic("unknown stop_on mode")
	}
}

func (gs *GameState) currentCell() *ui.Cell {
	return &gs.Cells[gs.CursorIndex]
}

func (gs *GameState) advance(correct bool) Action {
	cell := gs.currentCell()

	if correct {
		cell.State = ui.StatusCorrect
	} else {
		cell.State = ui.StatusIncorrect
		gs.Errors++
	}

	gs.CursorIndex++

	if gs.CursorIndex >= len(gs.Cells) {
		return gs.finish()
	}

	// we don't use cell because it got an outdated value
	if gs.currentCell().IsTerminator && !wordHasErrors(gs.Cells, gs.CursorIndex) {
		return gs.finish()
	}

	return ActionNone
}
