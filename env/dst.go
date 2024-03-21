package env

import (
	"log/slog"
	"strings"

	"github.com/yechentide/dstm/global"
	"github.com/yechentide/dstm/shell"
	"github.com/yechentide/dstm/utils"
)

func updateDSTServer(steamRootPath string, serverRootPath string, betaName string) error {
	scriptPath := steamRootPath + "/steamcmd.sh"
	args := []string{scriptPath, "+force_install_dir", serverRootPath, "+login", "anonymous", "+app_update", "343050"}
	if betaName != "" {
		slog.Debug("Installing beta: " + betaName)
		args = append(args, "-beta", betaName)
	}
	args = append(args, "validate", "+quit")

	cmd := strings.Join(args, " ")
	return shell.CreateTmuxSession(global.SESSION_DST_INSTALL, cmd)
}

func PrepareLatestDSTServer(steamRootPath string, serverRootPath string, betaName string) error {
	exists, err := IsSteamAvailable(steamRootPath)
	if err != nil {
		return err
	}
	if !exists {
		err = PrepareLatestSteam(steamRootPath)
		if err != nil {
			return err
		}
	}
	return updateDSTServer(steamRootPath, serverRootPath, betaName)
}

func IsDSTServerAvailable(serverRootPath string) (bool, error) {
	exists, err := utils.FileExists(serverRootPath + "/bin64/dontstarve_dedicated_server_nullrenderer_x64")
	if err != nil {
		return false, err
	}
	return exists, nil
}
