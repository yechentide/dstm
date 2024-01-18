/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yechentide/dstm/extractor"
	"github.com/yechentide/dstm/utils"
)

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract cluster settings in JSON",
	Long:  "Extract cluster settings in JSON format from dst server files",
	Run: func(cmd *cobra.Command, args []string) {
		serverRoot := viper.GetString("serverRoot")
		if serverRoot == "" {
			slog.Error("Please use --dst-root flag to specify dst root directory")
			os.Exit(1)
		}
		zipFile := utils.ExpandPath(serverRoot) + "/data/databundles/scripts.zip"
		tmpDir, err := cmd.Flags().GetString("output")
		if err != nil || tmpDir == "" {
			slog.Error("Please use --output flag to specify output directory")
			os.Exit(1)
		}
		extractor.ExtractSettings(zipFile, tmpDir)
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)

	extractCmd.Flags().StringP("output", "o", "", "output directory for cluster settings in JSON format")
}
