package main

import (
	"github.com/eduqr/ktype/internal/config"
	"github.com/eduqr/ktype/internal/loader"
	"github.com/eduqr/ktype/internal/logger"
	"github.com/eduqr/ktype/internal/typer"
	"github.com/eduqr/ktype/internal/ui"
	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) {
	cfg := config.Cfg
	log := logger.Setup(cfg.Logs.Enabled)

	log.Info("Starting ktype", "wordCount", cfg.Challenge.Words, "language", cfg.Challenge.Language, "quoteMode", cfg.Challenge.Quote)

	styles, err := ui.NewStyles(cfg.Theme)
	if err != nil {
		log.Error("Failed to setup styles", "error", err)
		return
	}

	challenge, err := loader.GenerateChallenge(cfg.Challenge)
	if err != nil {
		log.Error("Failed to generate challenge", "error", err)
		return
	}

	log.Info("Generated challenge", "count", len(challenge), "challenge", challenge)

	screen, err := ui.InitScreen()
	if err != nil {
		log.Error("Failed to initialize screen", "error", err)
		return
	}
	defer screen.Fini()

	w, h := screen.Size()
	cells := typer.BuildCells(challenge, 0, 1, w, h)

	ui.DrawChallenge(screen, cells, styles)

	gameState := typer.InitGameState(cells, styles)
	typer.StartGame(screen, gameState)
}
