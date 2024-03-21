package server

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/yechentide/dstm/logger"
	"github.com/yechentide/dstm/server"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart dstm server",
	Long:  "Restart dstm server.",
	Run: func(cmd *cobra.Command, args []string) {
		err := server.RestartShard(targetCluster, targetShard, skipModUpdate, forceShutdown)
		if err != nil {
			slog.Error("Something went wrong.", "error", err)
			logger.PrintJsonResultAndExit(1)
		}
	},
}

func init() {
	restartCmd.Flags().BoolVarP(&skipModUpdate, "skip-mod-update", "n", false, "skip mod update")
	restartCmd.Flags().BoolVarP(&forceShutdown, "force-shutdown", "f", false, "force shutdown")
}
