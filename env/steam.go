package env

import (
	"log/slog"
	"os"
	"strings"

	"github.com/yechentide/dstm/shell"
	"github.com/yechentide/dstm/utils"
	"golang.org/x/exp/slices"
)

func steamScriptExists(steamRootPath string) (bool, error) {
	exists, err := utils.FileExists(steamRootPath + "/steamcmd.sh")
	if err != nil {
		return false, err
	}
	return exists, nil
}

func downloadSteamScript(steamRootPath string) error {
	rootPath := steamRootPath
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

func updateSteam(steamRootPath string) error {
	scriptPath := steamRootPath + "/steamcmd.sh"
	args := []string{scriptPath, "+login", "anonymous", "validate", "+quit"}
	cmd := strings.Join(args, " ")
	return shell.CreateTmuxSession(TmuxSessionForSteam, cmd)
}

func PrepareLatestSteam(steamRootPath string) error {
	exists, err := steamScriptExists(steamRootPath)
	if err != nil {
		return err
	}
	if !exists {
		err = downloadSteamScript(steamRootPath)
		if err != nil {
			return err
		}
	}
	return updateSteam(steamRootPath)
}

func IsSteamAvailable(steamRootPath string) (bool, error) {
	exists, err := utils.FileExists(steamRootPath + "/linux32/steamclient.so")
	if err != nil {
		return false, err
	}
	return exists, nil
}
