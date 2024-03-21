package config

import (
	"github.com/spf13/cobra"
	"github.com/yechentide/dstm/logger"
	"github.com/yechentide/dstm/repl"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Config cluster/shard/world settings",
	Long:  "Config cluster.ini, shard.ini and worldgenoverride.lua.",
	Run: func(cmd *cobra.Command, args []string) {
		if isCluster {
			repl.UpdateCluster()
			logger.PrintJsonResultAndExit(0)
		}

		if isShard {
			repl.UpdateShard()
			logger.PrintJsonResultAndExit(0)
		}

		if isWorld {
			repl.UpdateWorld()
			logger.PrintJsonResultAndExit(0)
		}
	},
}

var (
	isCluster bool
	isShard   bool
	isWorld   bool
)

func init() {
	ConfigCmd.Flags().BoolVarP(&isCluster, "cluster", "c", false, "config cluster.ini")
	ConfigCmd.Flags().BoolVarP(&isShard, "shard", "s", false, "config shard.ini")
	ConfigCmd.Flags().BoolVarP(&isWorld, "world", "w", false, "config worldgenoverride.lua")
	ConfigCmd.MarkFlagsOneRequired("cluster", "shard", "world")
}
