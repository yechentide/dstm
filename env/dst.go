package env

import (
	"log/slog"
	"strings"

	"github.com/yechentide/dstm/shell"
	"github.com/yechentide/dstm/utils"
)

func updateDSTServer(steamRoot string, serverRoot string, betaName string) error {
	scriptPath := utils.ExpandPath(steamRoot + "/steamcmd.sh")
	args := []string{scriptPath, "+force_install_dir", serverRoot, "+login", "anonymous", "+app_update", "343050"}
	if betaName != "" {
		slog.Debug("Installing beta: " + betaName)
		args = append(args, "-beta", betaName)
	}
	args = append(args, "validate", "+quit")

	cmd := strings.Join(args, " ")
	return shell.CreateTmuxSession(TmuxSessionForDST, cmd)
}

func PrepareLatestDSTServer(steamRoot string, serverRoot string, betaName string) error {
	exists, err := IsSteamAvailable(steamRoot)
	if err != nil {
		return err
	}
	if !exists {
		err = PrepareLatestSteam(steamRoot)
		if err != nil {
			return err
		}
	}
	return updateDSTServer(steamRoot, serverRoot, betaName)
}

func IsDSTServerAvailable(serverRoot string) (bool, error) {
	exists, err := utils.FileExists(serverRoot + "/bin64/dontstarve_dedicated_server_nullrenderer_x64")
	if err != nil {
		return false, err
	}
	return exists, nil
}