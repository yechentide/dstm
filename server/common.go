package server

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
	"github.com/yechentide/dstm/shell"
)

func MakeSessionName(clusterName, shardName string) string {
	separator := viper.GetString("separator")
	return "dstm" + separator + clusterName + separator + shardName
}

func IsShardRunning(sessionName string) (bool, error) {
	exists, err := shell.HasTmuxSession(sessionName)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	_, err = GetPlayerNumber(sessionName)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func IsShardStarted(sessionName string) (bool, error) {
	log, err := shell.CaptureTmuxSessionOutput(sessionName, true)
	if err != nil {
		return false, err
	}
	lines := strings.Split(log, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if strings.HasSuffix(line, "Sim paused") {
			return true, nil
		}
		if strings.HasSuffix(line, "No auth token could be found.") {
			return false, errors.New("no available cluster token found")
		}
	}
	return false, nil
}
