package main

import (
	"github.com/eduqr/ktype/internal/config"
	"github.com/eduqr/ktype/internal/logger"
	"github.com/eduqr/ktype/internal/session"
	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) {
	cfg := config.Cfg
	log := logger.Setup(cfg.Logs.Enabled)

	session, err := session.New(cfg, log)
	if err != nil {
		log.Error("Failed to create session", "error", err)
		return
	}
	defer session.Close()

	session.Run()
}
