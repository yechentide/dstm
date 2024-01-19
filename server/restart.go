package server

func RestartShard(clusterName, shardName string, forceShutdown bool) error {
	err := StopShardIfExists(clusterName, shardName, forceShutdown)
	if err != nil {
		return err
	}
	return StartShard(clusterName, shardName, forceShutdown)
}
