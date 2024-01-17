package global

import (
	"log/slog"
	"os"
)

var logLevel = new(slog.LevelVar)

func SetDefaultLogger() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)
}

func UpdateLogLevel(level slog.Level) {
	logLevel.Set(level)
}
