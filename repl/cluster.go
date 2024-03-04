package repl

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/viper"
	"github.com/yechentide/dstm/config/cluster"
	"github.com/yechentide/dstm/utils"
)

func selectCluster(worldsDirPath string) string {
	existClusters, err := utils.ListAllClusters(worldsDirPath)
	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}
	selected := Selector(existClusters, "Please select a cluster", false)[0]
	return selected
}

func updateClusterConfig(cfg *cluster.ClusterConfig) {
	// TODO
	slog.Warn("Not implemented: UpdateClusterConfig()")
}

func showExistClusterList(worldsDirPath string) []string {
	existClusters, err := utils.ListAllClusters(worldsDirPath)
	if err == nil {
		msg := "Exist clusters:"
		for _, cluster := range existClusters {
			msg += " " + cluster
		}
		printInfo(msg)
	} else {
		existClusters = []string{}
		printError("Failed to list clusters")
	}
	fmt.Println()
	return existClusters
}

func CreateCluster() {
	worldsDirPath := viper.GetString("dataRootPath") + "/" + viper.GetString("worldsDirName")
	existClusters := showExistClusterList(worldsDirPath)

	clusterName := Readline("Enter cluster name (leave blank to cancel)", []func(string) error{
		utils.NotContainSpace,
		utils.Unique(existClusters),
	})
	if clusterName == "" {
		os.Exit(0)
	}

	clusterToken := Readline("Enter cluster token", []func(string) error{
		utils.NotContainSpace,
		utils.IsClusterToken,
	})

	newClusterPath := worldsDirPath + "/" + clusterName
	err := os.Mkdir(newClusterPath, 0755)
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
	updateClusterConfig(config)
	err = config.SaveTo(newClusterPath)
	if err != nil {
		printError("Failed to save cluster config: " + err.Error())
		os.Exit(1)
	}
}

func UpdateCluster() {
	worldsDirPath := viper.GetString("dataRootPath") + "/" + viper.GetString("worldsDirName")
	clusterDirPath := worldsDirPath + "/" + selectCluster(worldsDirPath)
	config, err := cluster.ReadClusterINI(clusterDirPath)
	if err != nil {
		printError("Failed to read cluster.ini in " + clusterDirPath)
		os.Exit(1)
	}
	updateClusterConfig(config)
	config.SaveTo(clusterDirPath)
}
