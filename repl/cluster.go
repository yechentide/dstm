package repl

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/yechentide/dstm/config/cluster"
	"github.com/yechentide/dstm/config/shard"
	"github.com/yechentide/dstm/utils"
)

func UpdateClusterConfig(cfg *cluster.ClusterConfig) {
	// TODO
	slog.Warn("Not implemented: UpdateClusterConfig()")
}

func CreateCluster() {
	worldsDir := viper.GetString("dataRootPath") + "/" + viper.GetString("worldsDirName")
	existClusters, err := utils.ListAllClusters(worldsDir)
	if err == nil {
		msg := "Exist clusters:"
		for _, cluster := range existClusters {
			msg += " " + cluster
		}
		printInfo(msg)
	} else {
		printError("Failed to list clusters")
	}
	fmt.Println()

	clusterName := Readline("Enter cluster name", []func(string) error{
		utils.NotEmpty,
		utils.NotContainSpace,
		utils.Unique(existClusters),
	})

	clusterToken := Readline("Enter cluster token", []func(string) error{
		utils.NotContainSpace,
		utils.IsClusterToken,
	})

	newClusterPath := worldsDir + "/" + clusterName
	err = utils.MkDirIfNotExists(newClusterPath, 0755, false)
	if err != nil {
		printError("Failed to create cluster: " + err.Error())
		os.Exit(1)
	}
	err = utils.WriteToFile(clusterToken, newClusterPath+"/cluster_token.txt")
	if err != nil {
		printError("Failed to write cluster token: " + err.Error())
		os.Exit(1)
	}

	config := cluster.MakeDefaultConfig()
	UpdateClusterConfig(config)
	err = config.SaveTo(newClusterPath)
	if err != nil {
		printError("Failed to save cluster config: " + err.Error())
		os.Exit(1)
	}
}

func UpdateShardConfig(config *shard.ShardConfig) {
	// TODO
	slog.Warn("Not implemented: UpdateShardConfig()")
}

func CreateShard(cluster string) {
	worldsDir := viper.GetString("dataRootPath") + "/" + viper.GetString("worldsDirName")
	clusterDir := worldsDir + "/" + cluster
	ok, err := utils.IsClusterDir(clusterDir)
	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}
	if !ok {
		printError("Not a valid cluster: " + cluster)
		os.Exit(1)
	}

	existShards, err := utils.ListShards(clusterDir)
	if err == nil {
		msg := "Exist shards:"
		for _, shard := range existShards {
			msg += " " + shard
		}
		printInfo(msg)
	} else {
		printError("can not list shards: " + err.Error())
		os.Exit(1)
	}
	printWarn("First shard will be generated as master shard.")
	shardType := Selector([]string{"forest", "cave"}, "Select shard type", false)[0]
	shardName := "Main"
	if len(existShards) > 0 {
		count := 0
		for _, shard := range existShards {
			name := strings.ToLower(shard)
			if strings.HasPrefix(name, shardType) {
				count++
			}
		}
		shardName = strings.ToUpper(string(shardType[0])) + shardType[1:] + fmt.Sprintf("%02d", count)
	}
	shardDir := clusterDir + "/" + shardName
	err = utils.MkDirIfNotExists(shardDir, 0755, false)
	if err != nil {
		printError("Failed to create shard: " + err.Error())
		os.Exit(1)
	}

	// server.ini
	config := shard.MakeDefaultConfig(shardType, len(existShards) == 0)
	UpdateShardConfig(config)
	err = config.SaveTo(shardDir)
	if err != nil {
		printError("Failed to save shard config: " + err.Error())
		os.Exit(1)
	}
	// worldgenoverride.lua
}
