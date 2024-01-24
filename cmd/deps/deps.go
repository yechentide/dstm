package deps

import (
	"github.com/spf13/cobra"
)

var DepsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Install dependencies",
	Long:  "Install dependencies.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	DepsCmd.AddCommand(listCmd)
	DepsCmd.AddCommand(installCmd)
}
