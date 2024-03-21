package logger

import (
	"context"
	"io"
	"log"
	"log/slog"

	"github.com/yechentide/dstm/global"
)

func newTextHandler(out io.Writer, level slog.Level) *textHandler {
	return &textHandler{
		Handler: slog.NewTextHandler(out, &slog.HandlerOptions{
			Level: level,
		}),
		logger: log.New(out, "", 0),
	}
}

type textHandler struct {
	slog.Handler
	logger *log.Logger
}

func (h *textHandler) Handle(_ context.Context, record slog.Record) error {
	timestamp := buildTimestamp(record)
	prefix := buildPrefix(record.Level)
	msg := buildMessage(record)

	if showColor {
		timestamp = global.GrayStyle.Render(timestamp)
		prefix = coloredLogLevel(record.Level, prefix, true)
		msg = coloredLogLevel(record.Level, msg, false)
	}

	h.logger.Println(timestamp + prefix + msg)
	return nil
}

func coloredLogLevel(level slog.Level, text string, isPrefix bool) string {
	switch level {
	case slog.LevelDebug:
		if isPrefix {
			return global.DebugStyle.Render(text)
		} else {
			return text
		}
	case slog.LevelWarn:
		return global.WarnStyle.Render(text)
	case slog.LevelInfo:
		if isPrefix {
			return global.InfoStyle.Render(text)
		} else {
			return text
		}
	case slog.LevelError:
		return global.ErrorStyle.Render(text)
	default:
		return text
	}
}
