package extractor

import (
	"errors"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/yechentide/dstm/global"
	"github.com/yechentide/dstm/shell"
	"github.com/yechentide/dstm/utils"
)

func ExtractWorldgenVanillaSettings(serverRootPath, outputDirPath string) error {
	zipFilePath := serverRootPath + "/data/databundles/scripts.zip"
	tmpDirPath := global.DIR_PATH_EXTRACT_WORKSPACE

	err := os.RemoveAll(tmpDirPath)
	if err != nil {
		return err
	}

	luaScriptsDirPath, err := prepareLuaScripts()
	if err != nil {
		return err
	}

	slog.Info("Unzipping scripts to " + tmpDirPath)
	err = utils.Unzip(zipFilePath, tmpDirPath, 0755)
	if err != nil {
		return err
	}

	unzippedScriptsDir := tmpDirPath + "/scripts"
	serverVersion := getServerVersion(serverRootPath)

	return executeLuaScriptToExtractSettings(luaScriptsDirPath, unzippedScriptsDir, outputDirPath, serverVersion)
}

func getServerVersion(serverRootPath string) string {
	versionFile := serverRootPath + "/version.txt"
	exists, err := utils.FileExists(versionFile)
	if err != nil || !exists {
		return "0"
	}
	data, err := os.ReadFile(versionFile)
	if err != nil {
		return "0"
	}
	return strings.TrimSpace(string(data))
}

func executeLuaScriptToExtractSettings(luaScriptsDirPath, unzippedScriptsDir, outputDirPath, serverVersion string) error {
	err := os.MkdirAll(outputDirPath, 0755)
	if err != nil {
		return err
	}

	slog.Info("Extracting worldgen vanilla settings ...")

	sessionName := "dstm-extract-worldgen-vanilla"
	cmd := "cd '" + luaScriptsDirPath + "' && lua extract-worldgen-vanilla.lua '" + serverVersion + "' '" + unzippedScriptsDir + "' '" + outputDirPath + "'"
	err = shell.CreateTmuxSession(sessionName, cmd)
	if err != nil {
		return err
	}
	for {
		time.Sleep(1 * time.Second)
		sessionExists, err := shell.HasTmuxSession(sessionName)
		if err != nil {
			return err
		}
		if !sessionExists {
			break
		}
	}
	exists, err := utils.FileExists(outputDirPath + "/en.forest.master.json")
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("failed to extract settings")
	}
	slog.Info("Worldgen vanilla settings extracted to " + outputDirPath)
	return nil
}
