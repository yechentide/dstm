package extract

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yechentide/dstm/extractor"
	"github.com/yechentide/dstm/utils"
)

var ExtractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract cluster settings in JSON",
	Long:  "Extract cluster settings in JSON format from dst server files",
	Run: func(cmd *cobra.Command, args []string) {
		serverRootPath := viper.GetString("serverRootPath")
		if serverRootPath == "" {
			slog.Error("Please use --server-root flag or config file to specify dst server root directory")
			os.Exit(1)
		}

		outputDirPath := viper.GetString("cacheDirPath") + "/json"
		specifiedDir, err := cmd.Flags().GetString("output")
		if err == nil && specifiedDir != "" {
			outputDirPath = utils.ExpandPath(specifiedDir)
		}

		extractor.ExtractSettings(serverRootPath, outputDirPath)
	},
}

func init() {
	ExtractCmd.Flags().StringP("output", "o", "", "output directory for cluster settings in JSON format")
}
