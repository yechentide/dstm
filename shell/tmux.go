package shell

import (
	"errors"
	"log/slog"
	"strings"
)

func ListTmuxSessions() ([]string, error) {
	output, err := ExecuteAndGetOutput("tmux", "list-sessions", "-F", "#S")
	if err == nil {
		sessions := strings.Split(output, "\n")
		return sessions, nil
	}
	if strings.HasPrefix(output, "no server running") {
		return []string{}, nil
	}
	return nil, errors.New(output)
}

func HasTmuxSession(sessionName string) (bool, error) {
	output, err := ExecuteAndGetOutput("tmux", "has-session", "-t", sessionName)
	if err == nil {
		return true, nil
	}
	if strings.HasPrefix(output, "no server running") ||
		strings.HasPrefix(output, "can't find session") ||
		strings.HasSuffix(output, "(No such file or directory)") {
		return false, nil
	}
	return false, errors.New(output)
}

func KillTmuxSession(sessionName string) error {
	output, err := ExecuteAndGetOutput("tmux", "kill-session", "-t", sessionName)
	if err != nil {
		return errors.New(output)
	}
	return nil
}

func CreateTmuxSession(sessionName string, cmd string) error {
	exists, err := HasTmuxSession(sessionName)
	if err != nil {
		return err
	}
	if exists {
		slog.Debug("Tmux session \"" + sessionName + "\" already exists")
		return nil
	}
	args := []string{"new-session", "-d"}
	if len(sessionName) > 0 {
		args = append(args, "-s", sessionName)
	}
	if len(cmd) > 0 {
		args = append(args, cmd)
	}
	output, err := ExecuteAndGetOutput("tmux", args...)
	if err != nil {
		return errors.New(output)
	}
	return nil
}

func SendMessageToTmuxSession(sessionName, message string, enter bool) error {
	var output string
	var err error
	if enter {
		output, err = ExecuteAndGetOutput("tmux", "send-keys", "-t", sessionName, message, "ENTER")
	} else {
		output, err = ExecuteAndGetOutput("tmux", "send-keys", "-t", sessionName, message)
	}
	if err != nil {
		return errors.New(output)
	}
	return nil
}

func CaptureTmuxSessionOutput(sessionName string, visibleOnly bool) (string, error) {
	args := []string{"capture-pane", "-p", "-t", sessionName}
	if visibleOnly {
		args = append(args, "-S", "-", "-E", "-")
	}
	output, err := ExecuteAndGetOutput("tmux", args...)
	if err != nil {
		return "", err
	}
	return output, nil
}

func LogTmuxSessionOutput(sessionName, logPath string) error {
	_, err := ExecuteAndGetOutput("tmux", "pipe-pane", "-t", sessionName, "-o", "cat > "+logPath)
	return err
}
