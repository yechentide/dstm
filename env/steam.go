package env

import (
	"log/slog"
	"os"
	"strings"

	"github.com/yechentide/dstm/shell"
	"github.com/yechentide/dstm/utils"
	"golang.org/x/exp/slices"
)

func steamScriptExists(steamRoot string) (bool, error) {
	exists, err := utils.FileExists(steamRoot + "/steamcmd.sh")
	if err != nil {
		return false, err
	}
	return exists, nil
}

func downloadSteamScript(steamRoot string) error {
	rootPath := utils.ExpandPath(steamRoot)
	err := utils.MkDirIfNotExists(rootPath, 0755, true)
	if err != nil {
		return err
	}

	dst := rootPath + "/steamcmd_linux.tar.gz"
	url := "https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz"
	err = utils.DownloadFile(dst, url)
	if err != nil {
		return err
	}

	output, err := shell.ExecuteAndGetOutput("tar", "-xvzf", dst, "--directory", rootPath)
	if err != nil {
		return err
	}
	extracted := strings.Split(output, "\n")
	containsScript := slices.Contains(extracted, "steamcmd.sh")
	if !containsScript {
		slog.Error("Failed to extract steamcmd.sh")
		os.Exit(1)
	}
	return nil
}

func updateSteam(steamRoot string) error {
	scriptPath := utils.ExpandPath(steamRoot + "/steamcmd.sh")
	args := []string{scriptPath, "+login", "anonymous", "validate", "+quit"}
	cmd := strings.Join(args, " ")
	return shell.CreateTmuxSession(TmuxSessionForSteam, cmd)
}

func PrepareLatestSteam(steamRoot string) error {
	exists, err := steamScriptExists(steamRoot)
	if err != nil {
		return err
	}
	if !exists {
		err = downloadSteamScript(steamRoot)
		if err != nil {
			return err
		}
	}
	return updateSteam(steamRoot)
}

func IsSteamAvailable(steamRoot string) (bool, error) {
	exists, err := utils.FileExists(steamRoot + "/linux32/steamclient.so")
	if err != nil {
		return false, err
	}
	return exists, nil
}
