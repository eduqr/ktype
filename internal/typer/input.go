package typer

import (
	"slices"
	"time"

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

func (gs *GameState) handleKey(screen tcell.Screen, ev *tcell.EventKey) Action {
	if ev.Key() == tcell.KeyCtrlC {
		return ActionQuit
	}

	if ev.Key() == tcell.KeyEscape {
		return ActionAbortRound
	}

	if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
		gs.handleBackspace(screen)
		return ActionNone
	}

	if ev.Key() == tcell.KeyRune {
		return gs.handleRune(screen, ev)
	}

	return ActionNone
}

func (gs *GameState) handleBackspace(screen tcell.Screen) {
	if gs.CursorIndex <= 0 || gs.IsDone {
		return
	}

	prevIndex := gs.CursorIndex - 1

	if gs.Cells[prevIndex].IsOverflow {
		gs.Cells = slices.Delete(gs.Cells, prevIndex, prevIndex+1)
		gs.CursorIndex--
		gs.OverflowCount--
		ui.RecalcPositions(gs.Cells, 0, 1, gs.screenWidth)
		ui.DrawChallenge(screen, gs.Cells, gs.Styles)
		return
	}

	// block crossing into the prev word if it was completed correctly
	if gs.Cells[gs.CursorIndex].WordIndex != gs.Cells[prevIndex].WordIndex {
		if isPrevWordCorrect(gs.Cells, gs.CursorIndex) {
			return
		}
	}

	ui.DrawCell(screen, gs.Cells[gs.CursorIndex], gs.Styles)
	gs.CursorIndex--
	gs.Cells[gs.CursorIndex].State = ui.StatusPending
}

const maxOverflow = 20

func (gs *GameState) handleRune(screen tcell.Screen, ev *tcell.EventKey) Action {
	if gs.IsDone {
		return ActionDone
	}

	input := ev.Str()
	currentIndex := gs.CursorIndex
	expected := gs.Cells[currentIndex].Expected

	if currentIndex == 0 && input != " " {
		gs.StartTime = time.Now()
		gs.startWPMSampler()
	}

	if expected == " " && input != " " {
		if gs.OverflowCount >= maxOverflow {
			return ActionNone
		}
		overFlowCell := ui.Cell{
			Expected:   input,
			State:      ui.StatusIncorrect,
			WordIndex:  gs.Cells[currentIndex].WordIndex,
			IsOverflow: true,
		}
		gs.Cells = slices.Insert(gs.Cells, currentIndex, overFlowCell)
		gs.OverflowCount++
		gs.Errors++
		ui.RecalcPositions(gs.Cells, 0, 1, gs.screenWidth)
		ui.DrawChallenge(screen, gs.Cells, gs.Styles)
		gs.CursorIndex++
		return ActionNone
	}

	if expected == " " && input == " " && gs.OverflowCount > 0 {
		return ActionNone
	}

	if input == " " && expected != " " {
		return ActionNone
	}

	if input == expected {
		gs.Cells[currentIndex].State = ui.StatusCorrect
	} else {
		gs.Cells[currentIndex].State = ui.StatusIncorrect
		gs.Errors++
	}

	ui.DrawCell(screen, gs.Cells[currentIndex], gs.Styles)
	gs.CursorIndex++

	if gs.CursorIndex >= len(gs.Cells) {
		gs.IsDone = true
		gs.EndTime = time.Now()
		return ActionDone
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

func WaitResultsAction(screen tcell.Screen) Action {
	for {
		event := <-screen.EventQ()
		if ev, ok := event.(*tcell.EventKey); ok {
			action := handleResultsKey(ev)
			if action != ActionNone {
				return action
			}
		}
	}
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
