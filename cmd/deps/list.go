package deps

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yechentide/dstm/env"
	"github.com/yechentide/dstm/logger"
)

var showStatus = false

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List required dependencies",
	Long:    "List required dependencies.",
	Run: func(cmd *cobra.Command, args []string) {
		helper, err := env.GetOSHelper()
		if err != nil {
			slog.Error("Failed to get os helper", "error", err)
			logger.PrintJsonResultAndExit(1)
		}

		pkgs := helper.Dependencies()
		if showStatus {
			status, err := helper.IsInstalled(pkgs)
			if err != nil {
				slog.Error(err.Error())
				logger.PrintJsonResultAndExit(1)
			}
			for pkg, isInstalled := range status {
				if isInstalled {
					fmt.Printf("ok:%s\n", pkg)
				} else {
					fmt.Printf("--:%s\n", pkg)
				}
			}
		} else {
			fmt.Println(strings.Join(pkgs, " "))
		}
	},
}

func init() {
	listCmd.Flags().BoolVarP(&showStatus, "status", "s", false, "show if installed or not")
}
