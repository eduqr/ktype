package session

import (
	"fmt"
	"log/slog"

	"github.com/eduqr/ktype/internal/config"
	"github.com/eduqr/ktype/internal/loader"
	"github.com/eduqr/ktype/internal/stats"
	"github.com/eduqr/ktype/internal/typer"
	"github.com/eduqr/ktype/internal/ui"
	"github.com/gdamore/tcell/v3"
)

type Session struct {
	screen    tcell.Screen
	gameState typer.GameState
	config    config.Config
}

func New(cfg config.Config, log *slog.Logger) (*Session, error) {
	if log != nil {
		log.Info("Starting ktype", "wordCount", cfg.Challenge.Words, "source", cfg.Challenge.Source, "quoteMode", cfg.Challenge.Quote)
	}

	styles, err := ui.NewStyles(cfg.Theme)
	if err != nil {
		return nil, fmt.Errorf("setup styles: %w", err)
	}

	challenge, err := loader.GenerateChallenge(cfg.Challenge)
	if err != nil {
		return nil, fmt.Errorf("generate challenge: %w", err)
	}
	if len(challenge) == 0 {
		return nil, fmt.Errorf("generate challenge: no words returned")
	}

	if log != nil {
		log.Info("Generated challenge", "count", len(challenge), "challenge", challenge)
	}

	screen, err := ui.InitScreen()
	if err != nil {
		return nil, fmt.Errorf("initialize screen: %w", err)
	}

	layout := ui.NewLayout(screen)
	cells := ui.BuildCells(challenge, layout)
	if len(cells) == 0 {
		screen.Fini()
		return nil, fmt.Errorf("build cells: no cells produced")
	}

	rules := typer.Rules{
		StopOn:          cfg.Challenge.StopOn,
		MaxWordOverflow: 20,
	}

	return &Session{
		screen:    screen,
		gameState: typer.InitGameState(challenge, cells, styles, layout, rules),
		config:    cfg,
	}, nil
}

func (session *Session) Run() error {
	for {
		result := typer.StartGame(session.screen, &session.gameState)

		if result.Action == typer.ActionQuit {
			return nil
		}

		var redraw func()

		if result.Action == typer.ActionAbortRound {
			redraw = func() {
				ui.ShowAbortRoundScreen(session.screen, session.gameState.Styles)
			}
		}

		if result.Action == typer.ActionDone {
			stats := stats.ComputeStats(result.CorrectCells, result.Errors, session.gameState.StartTime, session.gameState.EndTime)
			redraw = func() {
				ui.ShowFinalScreen(session.screen, session.gameState.Styles, stats, session.gameState.WPMSamples)
			}
		}

		redraw()
		nextAction := typer.WaitResultsAction(session.screen, redraw)

		switch nextAction {
		case typer.ActionRepeat:
			session.resetRound(session.gameState.Challenge)
		case typer.ActionNewChallenge:
			challenge, err := loader.GenerateChallenge(session.config.Challenge)
			if err != nil {
				return fmt.Errorf("generate new challenge: %w", err)
			}
			session.resetRound(challenge)
		case typer.ActionQuit:
			return nil
		}
	}
}

func (session *Session) Close() {
	if session == nil || session.screen == nil {
		return
	}

	session.screen.Fini()
}

func (session *Session) resetRound(challenge []string) {
	layout := ui.NewLayout(session.screen)
	cells := ui.BuildCells(challenge, layout)
	session.gameState = typer.InitGameState(challenge, cells, session.gameState.Styles, layout, session.gameState.Rules)
}
