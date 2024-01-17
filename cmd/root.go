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
	"github.com/yechentide/dstm/global"
)

var rootCmd = &cobra.Command{
	Use:     "dstm",
	Version: "v0.0.1",
	Short:   "Tools for Don't Starve Together Dedicated Server",
	Long:    "Tools for Don't Starve Together Dedicated Server.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		global.SetDefaultLogger()
		isDebug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			slog.Error("Failed to get debug flag", err)
			return
		}
		if isDebug {
			global.UpdateLogLevel(slog.LevelDebug)
		}
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

var (
	steamRoot = ""
	dstRoot   = ""
	betaName  = ""
)

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.SetVersionTemplate("(*•ᴗ•*) " + rootCmd.Use + " " + rootCmd.Version + "\n")

	rootCmd.PersistentFlags().Bool("debug", false, "print debug messages")

	rootCmd.PersistentFlags().StringVar(&steamRoot, "steam-root", "", "steam root directory")
	rootCmd.PersistentFlags().StringVar(&dstRoot, "dst-root", "", "dst server root directory")
	rootCmd.PersistentFlags().StringVar(&betaName, "beta", "", "name of dst server beta version")
}
