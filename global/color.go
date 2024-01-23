package global

import (
	"log/slog"

	"github.com/charmbracelet/lipgloss"
)

var (
	grayStyle  = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "244", Dark: "244"})
	debugStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "165", Dark: "165"})
	warnStyle  = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "228", Dark: "228"})
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "210", Dark: "210"})
)

func coloredLogLevel(level slog.Level, text string) string {
	switch level.String() {
	case "DEBUG":
		return debugStyle.Render(text)
	case "WARN":
		return warnStyle.Render(text)
	case "ERROR":
		return errorStyle.Render(text)
	default:
		return text
	}
}
