package typer

import (
	"github.com/eduqr/ktype/internal/ui"
	"github.com/gdamore/tcell/v3"
)

type GameState struct {
	Challenge     []string
	Cells         []ui.Cell
	Styles        ui.Styles
	CursorIndex   int
	TotalChars    int
	Errors        int
	IsDone        bool
	OverflowCount int
	screenWidth   int
}

type GameResult struct {
	Errors     int
	TotalChars int
	Action     Action
}

func InitGameState(challenge []string, cells []ui.Cell, styles ui.Styles, screenWidth int) GameState {
	return GameState{
		Challenge:   challenge,
		Cells:       cells,
		Styles:      styles,
		CursorIndex: 0,
		TotalChars:  len(cells),
		Errors:      0,
		IsDone:      false,
		screenWidth: screenWidth,
	}
}

func StartGame(screen tcell.Screen, gameState *GameState) GameResult {
	ui.DrawChallenge(screen, gameState.Cells, gameState.Styles)
	drawCursor(screen, gameState)
	screen.Show()

	for {
		event := <-screen.EventQ()

		switch ev := event.(type) {
		case *tcell.EventResize:
			screen.Sync()
			gameState.relayout(screen)
			ui.DrawChallenge(screen, gameState.Cells, gameState.Styles)
		case *tcell.EventKey:
			action := gameState.handleKey(screen, ev)
			if action != ActionNone {
				screen.Show()
				return GameResult{
					Errors:     gameState.Errors,
					TotalChars: gameState.TotalChars,
					Action:     action,
				}
			}
		}

		drawCursor(screen, gameState)
		screen.Show()
	}
}

func drawCursor(screen tcell.Screen, gameState *GameState) {
	if gameState.CursorIndex >= len(gameState.Cells) {
		return
	}

	cell := gameState.Cells[gameState.CursorIndex]
	ui.DrawCursorCell(screen, cell, gameState.Styles)
	screen.ShowCursor(cell.X, cell.Y)
}

func (gs *GameState) relayout(screen tcell.Screen) {
	screenWidth, _ := screen.Size()
	gs.screenWidth = screenWidth
	ui.RecalcPositions(gs.Cells, 0, 1, screenWidth)
}
