package extractor

import (
	"embed"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/yechentide/dstm/utils"
)

//go:embed scripts
var scriptsDir embed.FS

func prepareLuaScripts() (string, error) {
	destDirPath := "/tmp/dst-extractor"
	dirExists, err := utils.DirExists(destDirPath)
	if err != nil {
		return "", err
	}

	if dirExists {
		if isLatestScripts(destDirPath) {
			return destDirPath, nil
		} else {
			err = os.RemoveAll(destDirPath)
			if err != nil {
				return "", err
			}
		}
	}

	err = utils.CopyEmbeddedDir(scriptsDir, "scripts", destDirPath)
	if err != nil {
		return "", err
	}
	return destDirPath, nil
}

func getEmbeddedScriptVersion() string {
	file, err := scriptsDir.Open("scripts/version.txt")
	if err != nil {
		return "0"
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return "0"
	}

	return strings.TrimSpace(string(bytes))
}

func isLatestScripts(scriptsDirPath string) bool {
	slog.Debug("Checking if scripts are up to date ...")
	versionFile := scriptsDirPath + "/version.txt"
	data, err := os.ReadFile(versionFile)
	if err != nil {
		return false
	}
	currentVer := strings.TrimSpace(string(data))
	embeddedVer := getEmbeddedScriptVersion()
	slog.Debug("Current version: " + currentVer + ", Embedded version: " + embeddedVer)
	return currentVer == embeddedVer
}
