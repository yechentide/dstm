package global

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type simpleHandler struct {
	slog.Handler
	logger *log.Logger
}

func newSimpleHandler(out io.Writer, level slog.Level) *simpleHandler {
	prefix := ""
	h := &simpleHandler{
		Handler: slog.NewJSONHandler(out, &slog.HandlerOptions{
			Level: level,
		}),
		logger: log.New(out, prefix, 0),
	}
	return h
}

func (h *simpleHandler) Handle(_ context.Context, record slog.Record) error {
	ts := record.Time.Format("[2006/01/02 15:04:05]")
	levelStr := fmt.Sprintf("%7s", "["+record.Level.String()+"]")
	msg := record.Message
	record.Attrs(func(attr slog.Attr) bool {
		value := attr.Value.String()
		if strings.IndexByte(value, ' ') >= 0 {
			msg += fmt.Sprintf(" %s=\"%s\"", attr.Key, value)
		} else {
			msg += fmt.Sprintf(" %s=%s", attr.Key, value)
		}
		return true
	})
	if !viper.GetBool("noColor") {
		ts = GrayStyle.Render(ts)
		levelStr = coloredLogLevel(record.Level, levelStr)
		msg = coloredLogLevel(record.Level, msg)
	}
	h.logger.Println(ts, levelStr, msg)
	return nil
}

const levelNone slog.Level = 999

var logLevel = new(slog.LevelVar)

func InitCustomLogLevelAndFormat() {
	logger := slog.New(newSimpleHandler(os.Stdout, slog.LevelInfo))
	slog.SetDefault(logger)

	level := viper.GetString("logLevel")
	ChangeLogLevel(level)

	logger = slog.New(newSimpleHandler(os.Stdout, logLevel.Level()))
	slog.SetDefault(logger)
}

func ChangeLogLevel(level string) {
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
		slog.Error("Invalid log level: " + level)
		slog.Info("Current log level: " + level)
		return
	}
}
