package logger

import (
	"io"
	"log/slog"
	"os"
)

func Setup(logsEnable bool) *slog.Logger {
	var logOutput io.Writer

	if logsEnable {
		logOutput = os.Stderr
	} else {
		logOutput = io.Discard
	}

	return slog.New(slog.NewTextHandler(logOutput, nil))
}
