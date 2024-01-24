package server

import (
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage dstm server",
	Long:  "Start or stop dstm server with this command.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var (
	targetCluster string
	targetShard   string
	skipModUpdate = false
	forceShutdown = false
)

func init() {
	ServerCmd.AddCommand(startCmd)
	ServerCmd.AddCommand(stopCmd)
	ServerCmd.AddCommand(restartCmd)

	ServerCmd.PersistentFlags().StringVarP(&targetCluster, "cluster", "c", "", "target cluster")
	ServerCmd.PersistentFlags().StringVarP(&targetShard, "shard", "s", "", "target shard")
	ServerCmd.MarkFlagRequired("cluster")
}
