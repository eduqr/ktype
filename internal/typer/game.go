package typer

import (
	"time"

	"github.com/eduqr/ktype/internal/config"
	"github.com/eduqr/ktype/internal/stats"
	"github.com/eduqr/ktype/internal/ui"
	"github.com/gdamore/tcell/v3"
)

type GameState struct {
	Challenge       []string
	Cells           []ui.Cell
	Styles          ui.Styles
	CursorIndex     int
	Errors          int
	IsDone          bool
	WordMaxOverflow int
	layout          ui.Layout
	StartTime       time.Time
	EndTime         time.Time
	WPMSamples      []float64
	wpmTicker       *time.Ticker
}

type GameResult struct {
	Cells        []ui.Cell
	Errors       int
	CorrectCells int
	Action       Action
}

func InitGameState(challenge []string, cells []ui.Cell, styles ui.Styles, layout ui.Layout) GameState {
	return GameState{
		Challenge:       challenge,
		Cells:           cells,
		Styles:          styles,
		CursorIndex:     0,
		Errors:          0,
		IsDone:          false,
		WordMaxOverflow: 20,
		layout:          layout,
	}
}

func StartGame(screen tcell.Screen, gameState *GameState, challengeConfig config.ChallengeConfig) GameResult {
	drawGame(screen, gameState)
	events := screen.EventQ()

	for {
		var wpmTick <-chan time.Time
		if gameState.wpmTicker != nil {
			wpmTick = gameState.wpmTicker.C
		}

		select {
		case event := <-events:
			switch ev := event.(type) {
			case *tcell.EventResize:
				screen.Sync()
				gameState.relayout(screen)
				drawGame(screen, gameState)
			case *tcell.EventKey:
				action := gameState.handleKey(ev, challengeConfig)
				if action != ActionNone {
					gameState.stopWPMSampler()
					return GameResult{
						Cells:        gameState.Cells,
						Errors:       gameState.Errors,
						CorrectCells: gameState.countCorrectCells(),
						Action:       action,
					}
				}
				drawGame(screen, gameState)
			}
		case <-wpmTick:
			gameState.sampleWPM()
		}
	}
}

func drawGame(screen tcell.Screen, gameState *GameState) {
	ui.RecalcChallengePositions(gameState.Cells, gameState.layout)
	ui.DrawChallenge(screen, gameState.Cells, gameState.Styles)
	if gameState.layout.IsSupported() {
		drawCursor(screen, gameState)
	} else {
		screen.HideCursor()
	}
	screen.Show()
}

func (gs *GameState) startWPMSampler() {
	if gs.wpmTicker != nil {
		return
	}

	gs.wpmTicker = time.NewTicker(1 * time.Second)
}

func (gs *GameState) sampleWPM() {
	correctCells := gs.countCorrectCells()
	elapsed := time.Since(gs.StartTime).Seconds()
	sample := stats.CalculateWPM(correctCells, elapsed)
	gs.WPMSamples = append(gs.WPMSamples, sample)
}

func (gs *GameState) stopWPMSampler() {
	if gs.wpmTicker == nil {
		return
	}

	gs.wpmTicker.Stop()
	gs.wpmTicker = nil
}

func drawCursor(screen tcell.Screen, gameState *GameState) {
	cell := gameState.Cells[gameState.CursorIndex]
	screen.ShowCursor(cell.X, cell.Y)
}

func (gs *GameState) relayout(screen tcell.Screen) {
	gs.layout = ui.NewLayout(screen)
}

func (gs *GameState) countCorrectCells() int {
	count := 0
	for _, cell := range gs.Cells {
		if cell.State == ui.StatusCorrect && !cell.IsTerminator {
			count++
		}
	}
	return count
}
