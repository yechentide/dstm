package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yechentide/dstm/env"
	"github.com/yechentide/dstm/shell"
)

var depsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Install dependencies",
	Long:  "Install dependencies.",
	Run: func(cmd *cobra.Command, args []string) {
		getHelper := func() env.OSHelper {
			helper, err := env.GetOSHelper()
			if err != nil {
				slog.Error("Failed to get os helper", err)
				os.Exit(1)
			}
			return helper
		}
		installPkgs := func(packages []string) {
			helper := getHelper()
			var err error
			if len(packages) == 0 {
				err = helper.InstallAllRequired()
			} else {
				err = helper.InstallPackages(packages)
			}
			if err != nil {
				slog.Error("Failed to install required packages", err)
				os.Exit(1)
			}
		}
		if isList {
			helper := getHelper()
			fmt.Println(strings.Join(helper.Dependencies(), " "))
			os.Exit(0)
		}
		if isInstallPkg {
			installPkgs(args)
			os.Exit(0)
		}

		checkTmuxSession := func(sessionName string) bool {
			sessionExists, err := shell.HasTmuxSession(sessionName)
			if err != nil {
				slog.Error("Failed to check tmux session", err)
				os.Exit(1)
			}
			return sessionExists
		}
		waitForCompletion := func(sessionName string, checkFunc func() bool) bool {
			for {
				time.Sleep(1 * time.Second)
				sessionExists := checkTmuxSession(sessionName)
				if !sessionExists {
					break
				}
			}
			return checkFunc()
		}

		steamRoot := viper.GetString("steamRoot")
		serverRoot := viper.GetString("serverRoot")

		checkSteamRoot := func() {
			if steamRoot == "" {
				slog.Error("Please use --steam-root flag or config file to specify steam root directory")
				os.Exit(1)
			}
		}
		checkSteamAvailable := func(steamRoot string) func() bool {
			return func() bool {
				steamOK, err := env.IsSteamAvailable(steamRoot)
				if err != nil {
					slog.Error("Failed to check steam availability", err)
					os.Exit(1)
				}
				return steamOK
			}
		}
		prepareSteam := func() {
			sessionExists := checkTmuxSession(env.TmuxSessionForSteam)
			if sessionExists {
				slog.Info("Steam is installing")
				return
			}
			checkSteamRoot()
			err := env.PrepareLatestSteam(steamRoot)
			if err != nil {
				slog.Error("Failed to prepare steam", err)
				os.Exit(1)
			}
			steamOK := waitForCompletion(env.TmuxSessionForSteam, checkSteamAvailable(steamRoot))
			if !steamOK {
				slog.Error("Steam installation failed")
				os.Exit(1)
			}
		}
		if isInstallSteam {
			prepareSteam()
			os.Exit(0)
		}

		checkServerRoot := func() {
			if serverRoot == "" {
				slog.Error("Please use --server-root flag or config file to specify dst root directory")
				os.Exit(1)
			}
		}
		checkDSTAvailable := func(serverRoot string) func() bool {
			return func() bool {
				dstOK, err := env.IsDSTServerAvailable(serverRoot)
				if err != nil {
					slog.Error("Failed to check dst availability", err)
					os.Exit(1)
				}
				return dstOK
			}
		}
		prepareDSTServer := func() {
			sessionExists := checkTmuxSession(env.TmuxSessionForDST)
			if sessionExists {
				slog.Info("Steam is installing")
				return
			}
			checkSteamRoot()
			steamOK := checkSteamAvailable(steamRoot)()
			if !steamOK {
				slog.Error("Steam installation failed")
				os.Exit(1)
			}
			checkServerRoot()
			err := env.PrepareLatestDSTServer(steamRoot, serverRoot, "")
			if err != nil {
				slog.Error("Failed to prepare dst server", err)
				os.Exit(1)
			}
			dstOK := waitForCompletion(env.TmuxSessionForDST, checkDSTAvailable(serverRoot))
			if !dstOK {
				slog.Error("DST installation failed")
				os.Exit(1)
			}
		}
		if isInstallDST {
			prepareDSTServer()
			os.Exit(0)
		}

		if isInstallAll {
			installPkgs([]string{})
			prepareSteam()
			prepareDSTServer()
			os.Exit(0)
		}

		os.Exit(2)
	},
}

var (
	isList         bool
	isInstallAll   bool
	isInstallPkg   bool
	isInstallSteam bool
	isInstallDST   bool
)

func init() {
	rootCmd.AddCommand(depsCmd)

	depsCmd.Flags().BoolVarP(&isList, "list", "l", false, "list required packages")
	depsCmd.Flags().BoolVarP(&isInstallAll, "install", "i", false, "install packages & steam & dst server")
	depsCmd.Flags().BoolVar(&isInstallPkg, "install-pkg", false, "install specified packages, or all required packages if not specified")
	depsCmd.Flags().BoolVar(&isInstallSteam, "install-steam", false, "install or update steam")
	depsCmd.Flags().BoolVar(&isInstallDST, "install-dst", false, "install or update dst server")

	depsCmd.MarkFlagsOneRequired("list", "install", "install-pkg", "install-steam", "install-dst")
	depsCmd.MarkFlagsMutuallyExclusive("list", "install", "install-pkg", "install-steam", "install-dst")
}
