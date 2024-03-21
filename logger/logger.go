package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var logLevel = new(slog.LevelVar)

const levelNone slog.Level = 999

var (
	showTimestamp bool
	showPrefix    bool
	showColor     bool
	useJson       bool
)

func InitCustomLogLevelAndFormat() {
	useJson = viper.GetBool("json")
	showTimestamp = !viper.GetBool("hideTimestamp")
	showPrefix = !viper.GetBool("hidePrefix")
	showColor = !viper.GetBool("noColor")
	level := viper.GetString("logLevel")
	changeLogLevel(level)

	if useJson {
		logger := slog.New(newJsonHandler(os.Stdout, logLevel.Level()))
		slog.SetDefault(logger)
	} else {
		logger := slog.New(newTextHandler(os.Stdout, logLevel.Level()))
		slog.SetDefault(logger)
	}
}

func changeLogLevel(level string) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		logLevel.Set(slog.LevelDebug)
	case "INFO":
		logLevel.Set(slog.LevelInfo)
	case "WARN":
		logLevel.Set(slog.LevelWarn)
	case "ERROR":
		logLevel.Set(slog.LevelError)
	case "NONE":
		logLevel.Set(levelNone)
	default:
		logLevel.Set(slog.LevelInfo)
	}
}

func buildTimestamp(record slog.Record) string {
	if showTimestamp {
		return record.Time.Format("[2006/01/02 15:04:05]") + " "
	} else {
		return ""
	}
}

func getLevelLabel(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return "[DEBUG]"
	case slog.LevelWarn:
		return "[WARN]"
	case slog.LevelInfo:
		return "[INFO]"
	case slog.LevelError:
		return "[ERROR]"
	default:
		return ""
	}
}

func buildPrefix(level slog.Level) string {
	if !showPrefix {
		return ""
	}
	levelLabel := getLevelLabel(level)
	return fmt.Sprintf("%7s ", levelLabel)
}

func buildMessage(record slog.Record) string {
	msg := record.Message
	record.Attrs(func(attr slog.Attr) bool {
		value := attr.Value.String()
		if strings.IndexByte(value, ' ') >= 0 {
			msg += fmt.Sprintf(" %s:\"%s\"", attr.Key, value)
		} else {
			msg += fmt.Sprintf(" %s:%s", attr.Key, value)
		}
		return true
	})
	return msg
}
