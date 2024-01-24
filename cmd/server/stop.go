package server

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/yechentide/dstm/server"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop dstm server",
	Long:  "Stop dstm server.",
	Run: func(cmd *cobra.Command, args []string) {
		err := server.StopShardIfExists(targetCluster, targetShard, forceShutdown)
		if err != nil {
			slog.Error("Something went wrong.", "error", err)
		}
	},
}
