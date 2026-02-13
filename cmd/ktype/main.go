package main

import (
	"log/slog"
	"os"
)

func main() {
	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	logger := slog.New(jsonHandler)
	
	logger.Info("ktype cli logs")
}