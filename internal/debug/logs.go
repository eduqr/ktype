package debug

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
	
	jsonHandler := slog.NewJSONHandler(logOutput, nil)
	logger := slog.New(jsonHandler)
	return logger
}
