package global

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

var logLevel = new(slog.LevelVar)

func SetDefaultLogger() {
	level := viper.GetString("logLevel")
	UpdateLogLevel(level)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)
}

func UpdateLogLevel(level string) {
	switch level {
	case "debug":
		logLevel.Set(slog.LevelDebug)
	case "info":
		logLevel.Set(slog.LevelInfo)
	case "warn":
		logLevel.Set(slog.LevelWarn)
	case "error":
		logLevel.Set(slog.LevelError)
	default:
		slog.Error("Invalid log level: " + level)
		slog.Info("Current log level: " + level)
		return
	}
	slog.Info("Use log level: " + level)
}
