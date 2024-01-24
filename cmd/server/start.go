package server

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/yechentide/dstm/server"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start dstm server",
	Long:  "Start dstm server.",
	Run: func(cmd *cobra.Command, args []string) {
		err := server.StartShard(targetCluster, targetShard, skipModUpdate)
		if err != nil {
			slog.Error("Something went wrong.", "error", err)
		}
	},
}

func init() {
	startCmd.Flags().BoolVarP(&skipModUpdate, "skip-mod-update", "n", false, "skip mod update")
	startCmd.Flags().BoolVarP(&forceShutdown, "force-shutdown", "f", false, "force shutdown")
}
