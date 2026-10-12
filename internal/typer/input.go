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

	// a wrong letter stays under the cursor so we have to clear the cell first
	if stopOnMode == config.StopOnLetter && gs.Cells[gs.CursorIndex].State == ui.StatusIncorrect && !gs.Cells[gs.CursorIndex].IsOverflow && !isIndexAtWordStart(gs.Cells, gs.CursorIndex) {
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

	currentIndex := gs.CursorIndex
	expected := gs.Cells[currentIndex].Expected
	stopOnMode := gs.Rules.StopOn

	if gs.StartTime.IsZero() && currentIndex == 0 && input != " " {
		gs.StartTime = time.Now()
		gs.startWPMSampler()
	}

	if isIndexAtWordStart(gs.Cells, currentIndex) && input == " " {
		return ActionNone
	}

	if !isIndexAtWordStart(gs.Cells, currentIndex) && input == " " && stopOnMode == config.StopOnOff {
		gs.skipRemainingCells()
		if gs.CursorIndex >= len(gs.Cells) {
			return gs.finish()
		}
		return ActionNone
	}

	if input == " " && stopOnMode == config.StopOnLetter && gs.Cells[currentIndex].IsOverflow {
		// the overflow cell is replaced by the real space so both expect " "
		gs.Cells = slices.Delete(gs.Cells, currentIndex, currentIndex+1)
	}

	// so, we are at the end of a word, we print overflow when the input is not a space or there are errors in the word
	if expected == " " && (input != " " || currentWordHasErrors(gs.Cells, currentIndex)) {
		if overflowCountForWord(gs.Cells, currentIndex) >= gs.Rules.MaxWordOverflow && stopOnMode != config.StopOnLetter {
			return ActionNone
		}

		if input == " " {
			input = "_"
		}

		switch stopOnMode {
		case config.StopOnOff, config.StopOnWord:
			gs.appendOverflow(input)
			gs.CursorIndex++
		case config.StopOnLetter:
			gs.setOverflowCell(input)
		}
		return ActionNone
	}

	correct := input == expected

	if !correct && stopOnMode == config.StopOnLetter {
		// no index increment to match behaviour
		gs.Cells[currentIndex].State = ui.StatusIncorrect
		gs.Errors++
		return ActionNone
	}
	if correct {
		gs.Cells[currentIndex].State = ui.StatusCorrect
	} else {
		gs.Cells[currentIndex].State = ui.StatusIncorrect
		gs.Errors++
	}
	gs.CursorIndex++

	if gs.Cells[gs.CursorIndex].IsTerminator && !currentWordHasErrors(gs.Cells, currentIndex) {
		return gs.finish()
	}

	return ActionNone
}

func handleResultsKey(ev *tcell.EventKey) Action {
	if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
		return ActionQuit
	}

	if ev.Key() == tcell.KeyTab {
		return ActionNewChallenge
	}

	if ev.Key() == tcell.KeyRune {
		input := ev.Str()
		if input == "r" || input == "R" {
			return ActionRepeat
		}
	}

	return ActionNone
}

func WaitResultsAction(screen tcell.Screen, redraw func()) Action {
	for {
		event := <-screen.EventQ()

		switch ev := event.(type) {
		case *tcell.EventKey:
			action := handleResultsKey(ev)
			if action != ActionNone {
				return action
			}
		case *tcell.EventResize:
			screen.Sync()
			redraw()
		}
	}
}

func (gs *GameState) finish() Action {
	gs.IsDone = true
	gs.EndTime = time.Now()
	return ActionDone
}

func isPrevWordCorrect(cells []ui.Cell, cursorIndex int) bool {
	wordIndex := cells[cursorIndex-1].WordIndex
	for i := cursorIndex - 1; i >= 0 && cells[i].WordIndex == wordIndex; i-- {
		if cells[i].State != ui.StatusCorrect {
			return false
		}
	}
	return true
}

func currentWordHasErrors(cells []ui.Cell, cursorIndex int) bool {
	return wordHasState(cells, cursorIndex, ui.StatusIncorrect) || wordHasState(cells, cursorIndex, ui.StatusMissed)

}

func wordHasState(cells []ui.Cell, cursorIndex int, state ui.CellStatus) bool {
	wordIndex := cells[cursorIndex].WordIndex

	for i := wordStartIndex(cells, cursorIndex); i < len(cells) && cells[i].WordIndex == wordIndex; i++ {
		if cells[i].State == state {
			return true
		}
	}
	return false
}

func wordStartIndex(cells []ui.Cell, cursorIndex int) int {
	wordIndex := cells[cursorIndex].WordIndex

	start := cursorIndex
	for start > 0 && cells[start-1].WordIndex == wordIndex {
		start--
	}

	return start
}

func overflowCountForWord(cells []ui.Cell, cursorIndex int) int {
	wordIndex := cells[cursorIndex].WordIndex
	count := 0
	for i := wordStartIndex(cells, cursorIndex); i < len(cells) && cells[i].WordIndex == wordIndex; i++ {
		if cells[i].IsOverflow == true {
			count++
		}
	}
	return count
}

func isIndexAtWordStart(cells []ui.Cell, cursorIndex int) bool {
	return cursorIndex == 0 ||
		cells[cursorIndex-1].WordIndex != cells[cursorIndex].WordIndex
}

func (gs *GameState) skipRemainingCells() {
	wordIndex := gs.Cells[gs.CursorIndex].WordIndex
	skipped := false

	for gs.CursorIndex < len(gs.Cells) && gs.Cells[gs.CursorIndex].WordIndex == wordIndex {
		if gs.Cells[gs.CursorIndex].Expected == " " {
			gs.Cells[gs.CursorIndex].State = ui.StatusCorrect
		} else if gs.Cells[gs.CursorIndex].State == ui.StatusPending {
			gs.Cells[gs.CursorIndex].State = ui.StatusMissed
			skipped = true
		}
		gs.CursorIndex++
	}

	if skipped {
		gs.Errors++
	}
}

func (gs *GameState) appendOverflow(input string) {
	overflowCell := ui.Cell{
		Expected:     " ", // always space for overflow cells
		State:        ui.StatusIncorrect,
		WordIndex:    gs.Cells[gs.CursorIndex].WordIndex,
		OverflowChar: input,
		IsOverflow:   true,
	}
	gs.Cells = slices.Insert(gs.Cells, gs.CursorIndex, overflowCell)
	gs.Errors++
}

func (gs *GameState) setOverflowCell(input string) {
	if gs.Cells[gs.CursorIndex].IsOverflow {
		gs.Cells[gs.CursorIndex].OverflowChar = input
	} else {
		gs.appendOverflow(input)
	}
}
