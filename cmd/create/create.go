package create

import (
	"github.com/spf13/cobra"
	"github.com/yechentide/dstm/logger"
	"github.com/yechentide/dstm/repl"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create clusters and shards",
	Long:  "Create clusters and shards.",
	Run: func(cmd *cobra.Command, args []string) {
		if isCluster {
			repl.CreateCluster()
			logger.PrintJsonResultAndExit(0)
		}
		if isShard {
			repl.CreateShard()
			logger.PrintJsonResultAndExit(0)
		}
	},
}

var (
	isCluster bool
	isShard   bool
)

func init() {
	CreateCmd.Flags().BoolVarP(&isCluster, "cluster", "c", false, "create cluster")
	CreateCmd.Flags().BoolVarP(&isShard, "shard", "s", false, "create shard")
	CreateCmd.MarkFlagsOneRequired("cluster", "shard")
}
