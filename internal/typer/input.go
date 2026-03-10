package typer

import (
	"github.com/eduqr/ktype/internal/ui"
	"github.com/gdamore/tcell/v3"
)

type Action int

const (
	ActionNone Action = iota
	ActionQuit
	ActionDone
)

func (gs *GameState) handleKey(screen tcell.Screen, ev *tcell.EventKey) Action {
	if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
		return ActionQuit
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

	// block crossing into the prev word if it was completed correctly
	if gs.Cells[gs.CursorIndex].WordIndex != gs.Cells[gs.CursorIndex-1].WordIndex {
		if isPrevWordCorrect(gs.Cells, gs.CursorIndex) {
			return
		}
	}

	ui.DrawCell(screen, gs.Cells[gs.CursorIndex], gs.Styles)
	gs.CursorIndex--
	gs.Cells[gs.CursorIndex].State = ui.StatusPending
}

func (gs *GameState) handleRune(screen tcell.Screen, ev *tcell.EventKey) Action {
	if gs.IsDone {
		return ActionDone
	}

	input := ev.Str()
	currentIndex := gs.CursorIndex
	expected := gs.Cells[currentIndex].Expected

	if input == expected {
		gs.Cells[currentIndex].State = ui.StatusCorrect
	} else {
		gs.Cells[currentIndex].State = ui.StatusIncorrect
		gs.Errors++
	}

	ui.DrawCell(screen, gs.Cells[currentIndex], gs.Styles)
	gs.CursorIndex++

	if gs.CursorIndex >= gs.TotalChars {
		gs.IsDone = true
		return ActionDone
	}

	return ActionNone
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
