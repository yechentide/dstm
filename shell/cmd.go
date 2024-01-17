package shell

import (
	"log/slog"
	"os/exec"
	"strings"
)

func ExecuteAndGetOutput(cmd string, arg ...string) (string, error) {
	command := exec.Command(cmd, arg...)
	slog.Debug("ExecuteAndGetOutput()", "command", command.String())

	bytes, err := command.CombinedOutput()
	output := strings.TrimSpace(string(bytes))
	if err != nil {
		return output, err
	}
	return output, nil
}
