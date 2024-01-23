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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yechentide/dstm/global"
)

var rootCmd = &cobra.Command{
	Use:     "dstm",
	Version: "v0.0.1",
	Short:   "Tools for Don't Starve Together Dedicated Server",
	Long:    "Tools for Don't Starve Together Dedicated Server.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		global.InitCustomLogLevelAndFormat()
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

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.SetVersionTemplate("(*•ᴗ•*) " + rootCmd.Use + " " + rootCmd.Version + "\n")

	rootCmd.PersistentFlags().Bool("no-color", false, "disable color")
	viper.BindPFlag("noColor", rootCmd.PersistentFlags().Lookup("no-color"))

	rootCmd.PersistentFlags().String("log-level", "info", "change log error")
	viper.BindPFlag("logLevel", rootCmd.PersistentFlags().Lookup("log-level"))

	rootCmd.PersistentFlags().String("cache-dir", "$HOME/.cache/dstm", "cache directory")
	viper.BindPFlag("cacheDir", rootCmd.PersistentFlags().Lookup("cache-dir"))

	rootCmd.PersistentFlags().String("steam-root", "$HOME/Steam", "steam root directory")
	viper.BindPFlag("steamRoot", rootCmd.PersistentFlags().Lookup("steam-root"))

	rootCmd.PersistentFlags().String("server-root", "$HOME/DST/Server", "dst server root directory")
	viper.BindPFlag("serverRoot", rootCmd.PersistentFlags().Lookup("server-root"))

	rootCmd.PersistentFlags().String("data-root", "$HOME/DST/Klei", "dst save data root directory")
	viper.BindPFlag("dataRoot", rootCmd.PersistentFlags().Lookup("data-root"))

	rootCmd.PersistentFlags().String("separator", "-", "tmux session name separator")
	viper.BindPFlag("separator", rootCmd.PersistentFlags().Lookup("separator"))
}

func initConfig() {
	viper.SetDefault("noColor", false)
	viper.SetDefault("logLevel", "info")
	viper.SetDefault("cacheDir", "$HOME/.cache/dstm")
	viper.SetDefault("steamRoot", "$HOME/Steam")
	viper.SetDefault("serverRoot", "$HOME/DST/Server")
	viper.SetDefault("dataRoot", "$HOME/DST/Klei")
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
