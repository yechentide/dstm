package deps

import (
	"log/slog"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yechentide/dstm/env"
	"github.com/yechentide/dstm/global"
	"github.com/yechentide/dstm/logger"
	"github.com/yechentide/dstm/shell"
)

var helper env.OSHelper = nil

var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"i"},
	Short:   "Install dependencies",
	Long:    "Install dependencies.",
	Run: func(cmd *cobra.Command, args []string) {
		if helper == nil {
			h, err := env.GetOSHelper()
			if err != nil {
				slog.Error("Failed to get os helper.", "error", err)
				logger.PrintJsonResultAndExit(1)
			}
			helper = h
		}

		if pkgFlag {
			installPkgs(args)
			logger.PrintJsonResultAndExit(0)
		}

		isTerminalMultiplexerReady, err := helper.IsTerminalMultiplexerReady()
		if !isTerminalMultiplexerReady {
			slog.Error("Terminal Multiplexer is not available.")
			if err != nil {
				slog.Error(err.Error())
			}
			logger.PrintJsonResultAndExit(1)
		}

		steamRootPath := viper.GetString("steamRootPath")
		serverRootPath := viper.GetString("serverRootPath")

		if steamFlag {
			prepareSteam(steamRootPath)
			logger.PrintJsonResultAndExit(0)
		}

		if dstFlag {
			prepareDSTServer(steamRootPath, serverRootPath)
			logger.PrintJsonResultAndExit(0)
		}

		if allFlag {
			installPkgs([]string{})
			prepareSteam(steamRootPath)
			prepareDSTServer(steamRootPath, serverRootPath)
			logger.PrintJsonResultAndExit(0)
		}

		logger.PrintJsonResultAndExit(2)
	},
}

var (
	allFlag   bool
	pkgFlag   bool
	steamFlag bool
	dstFlag   bool
	password  = ""
)

func init() {
	installCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "install packages & steam & dst server")
	installCmd.Flags().BoolVarP(&pkgFlag, "pkg", "p", false, "install specified packages, or all required packages if no one specified")
	installCmd.Flags().BoolVarP(&steamFlag, "steam", "s", false, "install or update steam")
	installCmd.Flags().BoolVarP(&dstFlag, "dst", "d", false, "install or update dst server")

	installCmd.MarkFlagsOneRequired("all", "pkg", "steam", "dst")
	installCmd.MarkFlagsMutuallyExclusive("all", "pkg", "steam", "dst")

	installCmd.Flags().StringVar(&password, "password", "", "password to use sudo")
}

func installPkgs(packages []string) {
	var err error
	if len(packages) == 0 {
		err = helper.InstallAllRequired(password)
	} else {
		err = helper.InstallPackages(packages, password)
	}
	if err != nil {
		slog.Error("Failed to install required packages", "error", err)
		logger.PrintJsonResultAndExit(1)
	}
}

func checkTmuxSession(sessionName string) bool {
	sessionExists, err := shell.HasTmuxSession(sessionName)
	if err != nil {
		slog.Error("Failed to check tmux session", "error", err)
		logger.PrintJsonResultAndExit(1)
	}
	return sessionExists
}

func waitForCompletion(sessionName string, checkFunc func() bool) bool {
	for {
		time.Sleep(5 * time.Second)
		sessionExists := checkTmuxSession(sessionName)
		if !sessionExists {
			break
		}
	}
	return checkFunc()
}

/* ---------- ---------- ---------- ---------- ---------- ---------- */
// Steam

func checkSteamRoot(steamRootPath string) {
	if steamRootPath == "" {
		slog.Error("Please use --steam-root-path flag or config file to specify steam root directory")
		logger.PrintJsonResultAndExit(1)
	}
}

func prepareSteam(steamRootPath string) {
	sessionExists := checkTmuxSession(global.SESSION_STEAM_INSTALL)
	if sessionExists {
		slog.Info("Steam is installing")
		return
	}
	checkSteamRoot(steamRootPath)
	err := env.PrepareLatestSteam(steamRootPath)
	if err != nil {
		slog.Error("Failed to prepare steam", "error", err)
		logger.PrintJsonResultAndExit(1)
	}
	steamOK := waitForCompletion(global.SESSION_STEAM_INSTALL, checkSteamAvailable(steamRootPath))
	if !steamOK {
		slog.Error("Steam installation failed")
		logger.PrintJsonResultAndExit(1)
	}
}

func checkSteamAvailable(steamRootPath string) func() bool {
	return func() bool {
		steamOK, err := env.IsSteamAvailable(steamRootPath)
		if err != nil {
			slog.Error("Failed to check steam availability", "error", err)
			logger.PrintJsonResultAndExit(1)
		}
		return steamOK
	}
}

/* ---------- ---------- ---------- ---------- ---------- ---------- */
// DST Server

func checkServerRoot(serverRootPath string) {
	if serverRootPath == "" {
		slog.Error("Please use --server-root-path flag or config file to specify dst root directory")
		logger.PrintJsonResultAndExit(1)
	}
}

func checkDSTAvailable(serverRootPath string) func() bool {
	return func() bool {
		dstOK, err := env.IsDSTServerAvailable(serverRootPath)
		if err != nil {
			slog.Error("Failed to check dst availability", "error", err)
			logger.PrintJsonResultAndExit(1)
		}
		return dstOK
	}
}

func prepareDSTServer(steamRootPath, serverRootPath string) {
	sessionExists := checkTmuxSession(global.SESSION_DST_INSTALL)
	if sessionExists {
		slog.Info("Steam is installing")
		return
	}
	checkSteamRoot(steamRootPath)
	steamOK := checkSteamAvailable(steamRootPath)()
	if !steamOK {
		slog.Error("Steam installation failed")
		logger.PrintJsonResultAndExit(1)
	}
	checkServerRoot(serverRootPath)
	err := env.PrepareLatestDSTServer(steamRootPath, serverRootPath, "")
	if err != nil {
		slog.Error("Failed to prepare dst server", "error", err)
		logger.PrintJsonResultAndExit(1)
	}
	dstOK := waitForCompletion(global.SESSION_DST_INSTALL, checkDSTAvailable(serverRootPath))
	if !dstOK {
		slog.Error("DST installation failed")
		logger.PrintJsonResultAndExit(1)
	}
}
