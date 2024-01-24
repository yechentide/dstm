package deps

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yechentide/dstm/env"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List required dependencies",
	Long:    "List required dependencies.",
	Run: func(cmd *cobra.Command, args []string) {
		getHelper := func() env.OSHelper {
			helper, err := env.GetOSHelper()
			if err != nil {
				slog.Error("Failed to get os helper", err)
				os.Exit(1)
			}
			return helper
		}
		helper := getHelper()
		fmt.Println(strings.Join(helper.Dependencies(), " "))
	},
}
