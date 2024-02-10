package create

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yechentide/dstm/repl"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create clusters and shards",
	Long:  "Create clusters and shards.",
	Run: func(cmd *cobra.Command, args []string) {
		if isCluster {
			repl.CreateCluster()
			os.Exit(0)
		}
		if isShard {
			repl.CreateShard()
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
