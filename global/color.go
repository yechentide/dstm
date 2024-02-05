package global

import (
	"log/slog"

	"github.com/charmbracelet/lipgloss"
)

var (
	GrayStyle  = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "244", Dark: "244"})
	DebugStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "165", Dark: "165"})
	WarnStyle  = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "228", Dark: "228"})
	ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "210", Dark: "210"})
	InfoStyle  = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "86", Dark: "86"})
	OkStyle    = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "82", Dark: "82"})
)

func coloredLogLevel(level slog.Level, text string) string {
	switch level.String() {
	case "DEBUG":
		return DebugStyle.Render(text)
	case "WARN":
		return WarnStyle.Render(text)
	case "ERROR":
		return ErrorStyle.Render(text)
	default:
		return text
	}
}
