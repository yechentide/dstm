package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/yechentide/dstm/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage dstm server",
	Long:  "Start or stop dstm server with this command.",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if isStart {
			err = server.StartShard(targetCluster, targetShard, skipModUpdate)
		} else if isStop {
			err = server.StopShardIfExists(targetCluster, targetShard, forceShutdown)
		} else if isReStart {
			err = server.RestartShard(targetCluster, targetShard, forceShutdown)
		} else {
			slog.Info("Please specify start or stop or restart")
			return
		}
		if err != nil {
			slog.Error("Something went wrong.", "error", err)
		}
	},
}

var (
	isStart       bool
	isStop        bool
	isReStart     bool
	targetCluster string
	targetShard   string
	skipModUpdate = false
	forceShutdown = false
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().BoolVar(&isStart, "start", false, "start dstm server")
	serverCmd.Flags().BoolVar(&isStop, "stop", false, "stop dstm server")
	serverCmd.Flags().BoolVar(&isReStart, "restart", false, "restart dstm server")

	serverCmd.MarkFlagsOneRequired("start", "stop", "restart")
	serverCmd.MarkFlagsMutuallyExclusive("start", "stop", "restart")

	serverCmd.Flags().StringVarP(&targetCluster, "cluster", "c", "", "target cluster")
	serverCmd.Flags().StringVarP(&targetShard, "shard", "s", "", "target shard")

	serverCmd.MarkFlagRequired("cluster")
	serverCmd.MarkFlagRequired("shard")

	serverCmd.Flags().BoolVarP(&skipModUpdate, "skip-mod-update", "n", false, "skip mod update")
	serverCmd.Flags().BoolVarP(&forceShutdown, "force-shutdown", "f", false, "force shutdown")
}
