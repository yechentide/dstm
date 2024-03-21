package extract

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yechentide/dstm/extractor"
	"github.com/yechentide/dstm/logger"
	"github.com/yechentide/dstm/utils"
)

var ExtractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract cluster settings in JSON",
	Long:  "Extract cluster settings in JSON format from dst server files",
	Run: func(cmd *cobra.Command, args []string) {
		if outputDirPath == "" {
			outputDirPath = viper.GetString("cacheDirPath") + "/json"
		} else {
			outputDirPath = utils.ExpandPath(outputDirPath)
		}

		if worldgenVanilla {
			serverRootPath := viper.GetString("serverRootPath")
			if serverRootPath == "" {
				slog.Error("Please use --server-root flag or config file to specify dst server root directory")
				logger.PrintJsonResultAndExit(1)
			}
			err := extractor.ExtractWorldgenVanillaSettings(serverRootPath, outputDirPath)
			if err == nil {
				return
			} else {
				slog.Error("Failed to extract worldgen vanilla settings", "error", err)
				logger.PrintJsonResultAndExit(1)
			}
		}

		if worldgenOverride {
			if shardDirPath == "" {
				slog.Error("Please use --shard-dir-path flag or config file to specify shard directory")
				logger.PrintJsonResultAndExit(1)
			} else {
				shardDirPath = utils.ExpandPath(shardDirPath)
			}
			err := extractor.ExtractWorldgenOverride(shardDirPath, outputDirPath)
			if err == nil {
				return
			} else {
				slog.Error("Failed to extract worldgen override settings", "error", err)
				logger.PrintJsonResultAndExit(1)
			}
		}

		if modConfig {
			if modDirPath == "" {
				slog.Error("Please use --mod-dir-path flag or config file to specify mod directory")
				logger.PrintJsonResultAndExit(1)
			} else {
				modDirPath = utils.ExpandPath(modDirPath)
			}
			err := extractor.ExtractModConfiguration(modDirPath, outputDirPath, langCode)
			if err == nil {
				return
			} else {
				slog.Error("Failed to extract mod config settings", "error", err)
				logger.PrintJsonResultAndExit(1)
			}
		}

		if modOverride {
			if shardDirPath == "" {
				slog.Error("Please use --shard-dir-path flag or config file to specify shard directory")
				logger.PrintJsonResultAndExit(1)
			} else {
				shardDirPath = utils.ExpandPath(shardDirPath)
			}
			err := extractor.ExtractModOverride(shardDirPath, outputDirPath)
			if err == nil {
				return
			} else {
				slog.Error("Failed to extract mod override settings", "error", err)
				logger.PrintJsonResultAndExit(1)
			}
		}

		slog.Warn("Nothing happened. Please use --help flag to show usage")
	},
}

var (
	worldgenVanilla  bool
	worldgenOverride bool
	modConfig        bool
	modOverride      bool

	langCode      string
	modDirPath    string
	shardDirPath  string
	outputDirPath string
)

func init() {
	ExtractCmd.Flags().BoolVar(&worldgenVanilla, "worldgen-vanilla", false, "extract worldgen vanilla settings")
	ExtractCmd.Flags().BoolVar(&worldgenOverride, "worldgen-override", false, "extract worldgen override settings")
	ExtractCmd.Flags().BoolVar(&modConfig, "mod-config", false, "extract mod config settings")
	ExtractCmd.Flags().BoolVar(&modOverride, "mod-override", false, "extract mod override settings")

	ExtractCmd.MarkFlagsOneRequired("worldgen-vanilla", "worldgen-override", "mod-config", "mod-override")
	ExtractCmd.MarkFlagsMutuallyExclusive("worldgen-vanilla", "worldgen-override", "mod-config", "mod-override")

	ExtractCmd.Flags().StringVarP(&langCode, "lang", "l", "en", "language code")
	ExtractCmd.Flags().StringVarP(&modDirPath, "mod-dir-path", "m", "", "mod directory")
	ExtractCmd.Flags().StringVarP(&shardDirPath, "shard-dir-path", "s", "", "shard directory")
	ExtractCmd.Flags().StringVarP(&outputDirPath, "output-dir-path", "o", "", "output directory for cluster settings in JSON format")
}
