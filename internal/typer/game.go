package typer

import (
	"time"

	"github.com/eduqr/ktype/internal/ui"
	"github.com/gdamore/tcell/v3"
)

type GameState struct {
	Challenge   []string
	Cells       []ui.Cell
	Styles      ui.Styles
	Rules       Rules
	CursorIndex int
	Errors      int
	IsDone      bool
	layout      ui.Layout
	StartTime   time.Time
	EndTime     time.Time
	WPMSamples  []float64
	wpmTicker   *time.Ticker
}

type GameResult struct {
	Cells        []ui.Cell
	Errors       int
	CorrectCells int
	Action       Action
}

func InitGameState(challenge []string, cells []ui.Cell, styles ui.Styles, layout ui.Layout, rules Rules) GameState {
	return GameState{
		Challenge:   challenge,
		Cells:       cells,
		Styles:      styles,
		CursorIndex: 0,
		Errors:      0,
		IsDone:      false,
		Rules:       rules,
		layout:      layout,
	}
}

func StartGame(screen tcell.Screen, gameState *GameState) GameResult {
	ui.DrawGame(screen, gameState.Cells, gameState.Styles, gameState.layout, gameState.CursorIndex)
	events := screen.EventQ()
	pasting := false

	for {
		var wpmTick <-chan time.Time
		if gameState.wpmTicker != nil {
			wpmTick = gameState.wpmTicker.C
		}

		select {
		case event := <-events:
			switch ev := event.(type) {
			case *tcell.EventPaste:
				pasting = ev.Start()
			case *tcell.EventResize:
				screen.Sync()
				gameState.relayout(screen)
				ui.DrawGame(screen, gameState.Cells, gameState.Styles, gameState.layout, gameState.CursorIndex)
			case *tcell.EventKey:
				if pasting {
					continue
				}
				action := gameState.handleKey(ev)
				if action != ActionNone {
					gameState.stopWPMSampler()
					return GameResult{
						Cells:        gameState.Cells,
						Errors:       gameState.Errors,
						CorrectCells: gameState.countCorrectCells(),
						Action:       action,
					}
				}
				ui.DrawGame(screen, gameState.Cells, gameState.Styles, gameState.layout, gameState.CursorIndex)
			}
		case <-wpmTick:
			gameState.sampleWPM()
		}
	}
}

func (gs *GameState) finish() Action {
	gs.IsDone = true
	gs.EndTime = time.Now()
	return ActionDone
}

func (gs *GameState) relayout(screen tcell.Screen) {
	gs.layout = ui.NewLayout(screen)
}
