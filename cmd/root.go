/*
Copyright © 2024 yechentide

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yechentide/dstm/cmd/deps"
	"github.com/yechentide/dstm/cmd/extract"
	"github.com/yechentide/dstm/cmd/server"
	"github.com/yechentide/dstm/global"
	"github.com/yechentide/dstm/utils"
)

var rootCmd = &cobra.Command{
	Use:     "dstm",
	Version: "v0.0.1",
	Short:   "Tools for Don't Starve Together Dedicated Server",
	Long:    "Tools for Don't Starve Together Dedicated Server.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		global.InitCustomLogLevelAndFormat()
		expandPaths()
		debugConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Flow: rootCmd.Execute --> flags processing --> cobra.OnInitialize --> rootCmd.Run
	cobra.OnInitialize(initConfig)
	addFlags()
	addCommands()
}

func addFlags() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.SetVersionTemplate("(*•ᴗ•*) " + rootCmd.Use + " " + rootCmd.Version + "\n")

	rootCmd.PersistentFlags().Bool("no-color", false, "disable color")
	viper.BindPFlag("noColor", rootCmd.PersistentFlags().Lookup("no-color"))

	rootCmd.PersistentFlags().String("log-level", "info", "change log error")
	viper.BindPFlag("logLevel", rootCmd.PersistentFlags().Lookup("log-level"))

	rootCmd.PersistentFlags().String("cache-dir-path", "$HOME/.cache/dstm", "path of cache directory")
	viper.BindPFlag("cacheDirPath", rootCmd.PersistentFlags().Lookup("cache-dir-path"))

	rootCmd.PersistentFlags().String("steam-root-path", "$HOME/Steam", "path of steam root directory")
	viper.BindPFlag("steamRootPath", rootCmd.PersistentFlags().Lookup("steam-root-path"))

	rootCmd.PersistentFlags().String("server-root-path", "$HOME/DST/Server", "path of dst server root directory")
	viper.BindPFlag("serverRootPath", rootCmd.PersistentFlags().Lookup("server-root-path"))

	rootCmd.PersistentFlags().String("data-root-path", "$HOME/DST/Klei", "path of the dst save data root directory")
	viper.BindPFlag("dataRootPath", rootCmd.PersistentFlags().Lookup("data-root-path"))

	rootCmd.PersistentFlags().String("worlds-dir", "Worlds", "worlds directory name")
	viper.BindPFlag("worldsDir", rootCmd.PersistentFlags().Lookup("worlds-dir"))

	rootCmd.PersistentFlags().String("separator", "-", "tmux session name separator")
	viper.BindPFlag("separator", rootCmd.PersistentFlags().Lookup("separator"))
}

func addCommands() {
	rootCmd.AddCommand(deps.DepsCmd)
	rootCmd.AddCommand(extract.ExtractCmd)
	rootCmd.AddCommand(server.ServerCmd)
}

func initConfig() {
	viper.SetDefault("noColor", false)
	viper.SetDefault("logLevel", "info")
	viper.SetDefault("cacheDirPath", "$HOME/.cache/dstm")
	viper.SetDefault("steamRootPath", "$HOME/Steam")
	viper.SetDefault("serverRootPath", "$HOME/DST/Server")
	viper.SetDefault("dataRootPath", "$HOME/DST/Klei")
	viper.SetDefault("worldsDirName", "Worlds")
	viper.SetDefault("separator", "-")

	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	// search paths
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		viper.AddConfigPath(xdgConfigHome + "/dstm")
	}
	viper.AddConfigPath("$HOME/.dstm")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			slog.Warn("Config file not found")
		} else {
			// Config file was found but another error was produced
			slog.Error("Failed to read config file", "error", err)
		}
	}
}

func expandPaths() {
	viper.Set("cacheDirPath", utils.ExpandPath(viper.GetString("cacheDirPath")))
	viper.Set("steamRootPath", utils.ExpandPath(viper.GetString("steamRootPath")))
	viper.Set("serverRootPath", utils.ExpandPath(viper.GetString("serverRootPath")))
	viper.Set("dataRootPath", utils.ExpandPath(viper.GetString("dataRootPath")))
}

func debugConfig() {
	slog.Debug("========== ========== ========== =========")
	slog.Debug("noColor: " + strconv.FormatBool(viper.GetBool("noColor")))
	slog.Debug("logLevel: " + viper.GetString("logLevel"))
	slog.Debug("cacheDirPath: " + viper.GetString("cacheDirPath"))
	slog.Debug("steamRootPath: " + viper.GetString("steamRootPath"))
	slog.Debug("serverRootPath: " + viper.GetString("serverRootPath"))
	slog.Debug("dataRootPath: " + viper.GetString("dataRootPath"))
	slog.Debug("worldsDirName: " + viper.GetString("worldsDirName"))
	slog.Debug("separator: " + viper.GetString("separator"))
	slog.Debug("========== ========== ========== =========")
}
