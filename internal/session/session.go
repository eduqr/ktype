package session

import (
	"fmt"
	"log/slog"

	"github.com/eduqr/ktype/internal/config"
	"github.com/eduqr/ktype/internal/loader"
	"github.com/eduqr/ktype/internal/typer"
	"github.com/eduqr/ktype/internal/ui"
	"github.com/gdamore/tcell/v3"
)

type Session struct {
	screen    tcell.Screen
	gameState typer.GameState
}

func New(cfg config.Config, log *slog.Logger) (*Session, error) {
	if log != nil {
		log.Info("Starting ktype", "wordCount", cfg.Challenge.Words, "language", cfg.Challenge.Language, "quoteMode", cfg.Challenge.Quote)
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

	screenWidth, _ := screen.Size()
	cells := ui.BuildCells(challenge, 0, 1, screenWidth)
	if len(cells) == 0 {
		screen.Fini()
		return nil, fmt.Errorf("build cells: no cells produced")
	}

	return &Session{
		screen:    screen,
		gameState: typer.InitGameState(challenge, cells, styles),
	}, nil
}

func (session *Session) Run() {
	typer.StartGame(session.screen, &session.gameState)
}

func (session *Session) Close() {
	if session == nil || session.screen == nil {
		return
	}

	session.screen.Fini()
}
