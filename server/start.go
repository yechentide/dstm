package server

import (
	"errors"
	"log/slog"
	"time"

	"github.com/spf13/viper"
	"github.com/yechentide/dstm/env"
	"github.com/yechentide/dstm/shell"
)

func StartShard(clusterName, shardName string, skipModUpdate bool) error {
	sessionName := MakeSessionName(clusterName, shardName)

	exists, err := IsShardRunning(sessionName)
	if err != nil {
		return err
	}
	if exists {
		slog.Info("Shard already running")
		return nil
	}

	serverRoot := viper.GetString("serverRoot")
	ugcDir := serverRoot + "/ugc_mods"
	dataRoot := viper.GetString("dataRoot")
	worldsDirName := "worlds"

	cmd := "cd " + serverRoot
	if env.Is64BitCPU() {
		cmd = cmd + "/bin64; ./dontstarve_dedicated_server_nullrenderer_x64"
	} else {
		cmd = cmd + "/bin; ./dontstarve_dedicated_server_nullrenderer"
	}
	if skipModUpdate {
		cmd += " -skip_update_server_mods"
	}
	cmd += " -ugc_directory " + ugcDir
	cmd += " -persistent_storage_root " + dataRoot
	cmd += " -conf_dir " + worldsDirName
	cmd += " -cluster " + clusterName
	cmd += " -shard " + shardName

	slog.Info("Starting shard: " + sessionName)
	err = shell.CreateTmuxSession(sessionName, cmd)
	if err != nil {
		_ = StopShardIfExists(clusterName, shardName, true)
		return err
	}

	for i := 0; i < 120; i++ {
		time.Sleep(1 * time.Second)
		ok, err := IsShardStarted(sessionName)
		if err != nil {
			_ = StopShardIfExists(clusterName, shardName, true)
			return err
		}
		if ok {
			slog.Info("Shard started: " + sessionName)
			return nil
		}
	}
	return errors.New("Cannot start shard: " + sessionName)
}
