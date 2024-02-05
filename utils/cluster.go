package utils

import "log/slog"

func IsClusterDir(dirPath string) (bool, error) {
	exists, err := DirExists(dirPath)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	return FileExists(dirPath + "/cluster.ini")
}

func IsShardDir(dirPath string) (bool, error) {
	exists, err := DirExists(dirPath)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	return FileExists(dirPath + "/server.ini")
}

func ListAllClusters(worldsDir string) ([]string, error) {
	worldsDir = ExpandPath(worldsDir)
	dirs, err := ListDirs(worldsDir)
	if err != nil {
		return nil, err
	}

	var clusters []string
	for _, dirName := range dirs {
		path := worldsDir + "/" + dirName
		isCluster, err := IsClusterDir(path)
		if err != nil {
			slog.Warn("Something went wrong.", "error", err)
			continue
		}
		if isCluster {
			clusters = append(clusters, dirName)
		}
	}
	return clusters, nil
}

func ListShards(clusterDir string) ([]string, error) {
	clusterDir = ExpandPath(clusterDir)
	dirs, err := ListDirs(clusterDir)
	if err != nil {
		return nil, err
	}

	var shards []string
	for _, dirName := range dirs {
		path := clusterDir + "/" + dirName
		isShard, err := IsShardDir(path)
		if err != nil {
			slog.Warn("Something went wrong.", "error", err)
			continue
		}
		if isShard {
			shards = append(shards, dirName)
		}
	}
	return shards, nil
}
