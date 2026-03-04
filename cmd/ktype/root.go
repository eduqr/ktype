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
	log := logger.Setup(config.LogsEnable)

	log.Info("Starting ktype", "wordCount", config.WordCount, "language", config.Language, "quoteMode", config.QuoteMode)

	challenge, err := loader.GenerateChallenge(config.Language, config.WordCount, config.QuoteMode)
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

	ui.DrawChallenge(screen, cells, ui.SetupStyles())

	gameState := typer.InitGameState(cells)
	typer.StartGame(screen, gameState)
}
