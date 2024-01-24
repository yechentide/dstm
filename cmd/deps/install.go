package deps

import (
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yechentide/dstm/env"
	"github.com/yechentide/dstm/shell"
)

var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"i"},
	Short:   "Install dependencies",
	Long:    "Install dependencies.",
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

		if pkgFlag {
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
		if steamFlag {
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
		if dstFlag {
			prepareDSTServer()
			os.Exit(0)
		}

		if allFlag {
			installPkgs([]string{})
			prepareSteam()
			prepareDSTServer()
			os.Exit(0)
		}

		os.Exit(2)
	},
}

var (
	allFlag   bool
	pkgFlag   bool
	steamFlag bool
	dstFlag   bool
)

func init() {
	installCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "install packages & steam & dst server")
	installCmd.Flags().BoolVarP(&pkgFlag, "pkg", "p", false, "install specified packages, or all required packages if no one specified")
	installCmd.Flags().BoolVarP(&steamFlag, "steam", "s", false, "install or update steam")
	installCmd.Flags().BoolVarP(&dstFlag, "dst", "d", false, "install or update dst server")

	installCmd.MarkFlagsOneRequired("all", "pkg", "steam", "dst")
	installCmd.MarkFlagsMutuallyExclusive("all", "pkg", "steam", "dst")
}
