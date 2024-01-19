package server

import (
	"errors"
	"log/slog"
	"time"

	"github.com/yechentide/dstm/shell"
)

func StopShardIfExists(clusterName, shardName string, forceShutdown bool) error {
	sessionName := MakeSessionName(clusterName, shardName)
	exists, err := shell.HasTmuxSession(sessionName)
	if err != nil {
		return err
	}
	if !exists {
		slog.Info("No shard running: " + sessionName)
		return nil
	}

	slog.Info("Stopping shard: " + sessionName)
	isRunning, err := IsShardRunning(sessionName)
	if err != nil {
		return err
	}
	if !isRunning {
		slog.Info("Shard already stopped: " + sessionName)
		return shell.KillTmuxSession(sessionName)
	}

	err = shell.SendMessageToTmuxSession(sessionName, "c_shutdown(true)", true)
	if err != nil {
		return err
	}

	for i := 0; i < 20; i++ {
		time.Sleep(1 * time.Second)
		exists, err := shell.HasTmuxSession(sessionName)
		if err != nil {
			return err
		}
		if !exists {
			slog.Info("Shard stopped: " + sessionName)
			return nil
		}
	}

	if forceShutdown {
		slog.Info("Force shutdown shard: " + sessionName)
		return shell.KillTmuxSession(sessionName)
	}
	return errors.New("Cannot stop shard: " + sessionName)
}
