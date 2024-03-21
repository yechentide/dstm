package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
)

var result *resultOutput = nil

type resultOutput struct {
	StatusCode int      `json:"status"`
	DebugMsgs  []string `json:"debug-messages"`
	WarnMsgs   []string `json:"warn-messages"`
	InfoMsgs   []string `json:"info-messages"`
	ErrMsgs    []string `json:"error-messages"`
}

func (r *resultOutput) reset() {
	r.StatusCode = -1
	r.DebugMsgs = []string{}
	r.WarnMsgs = []string{}
	r.InfoMsgs = []string{}
	r.ErrMsgs = []string{}
}

func PrintJsonResultAndExit(statusCode int) {
	PrintJsonResult(statusCode)
	os.Exit(statusCode)
}

func PrintJsonResult(statusCode int) {
	if result == nil {
		return
	}
	result.StatusCode = statusCode
	jsonData, err := json.Marshal(result)
	if err != nil {
		fmt.Println("{ \"status\": 1 }")
		os.Exit(1)
	}
	fmt.Println(string(jsonData))
	result.reset()
}

func newJsonHandler(out io.Writer, level slog.Level) *jsonHandler {
	result = &resultOutput{}
	result.reset()
	return &jsonHandler{
		Handler: slog.NewTextHandler(out, &slog.HandlerOptions{
			Level: level,
		}),
	}
}

type jsonHandler struct {
	slog.Handler
}

func (h *jsonHandler) Handle(_ context.Context, record slog.Record) error {
	timestamp := buildTimestamp(record)
	prefix := buildPrefix(record.Level)
	msg := buildMessage(record)

	msg = timestamp + prefix + msg

	switch record.Level {
	case slog.LevelDebug:
		result.DebugMsgs = append(result.DebugMsgs, msg)
	case slog.LevelWarn:
		result.WarnMsgs = append(result.WarnMsgs, msg)
	case slog.LevelInfo:
		result.InfoMsgs = append(result.InfoMsgs, msg)
	case slog.LevelError:
		result.ErrMsgs = append(result.ErrMsgs, msg)
	default:
		return nil
	}
	return nil
}
