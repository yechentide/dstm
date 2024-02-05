package extractor

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yechentide/dstm/shell"
	"github.com/yechentide/dstm/utils"
)

//go:embed worldgen_template
var scriptsDir embed.FS

func ExtractSettings(serverRootPath, outputDirPath string) error {
	serverRootPath = utils.ExpandPath(serverRootPath)
	serverVersion := getServerVersion(serverRootPath)

	zipFile := serverRootPath + "/data/databundles/scripts.zip"
	zipFilePath := utils.ExpandPath(zipFile)
	tmpDir := "/tmp/dstm-extract-worldgen-template"
	scriptDir := tmpDir + "/scripts"
	workDir := tmpDir + "/work"

	err := utils.DelDirIfExists(tmpDir)
	if err != nil {
		return err
	}

	err = prepareFiles(zipFilePath, tmpDir, scriptDir, workDir)
	if err != nil {
		return err
	}

	outputDirPath = utils.ExpandPath(outputDirPath)
	return executeLuaScriptToExtractSettings(workDir, outputDirPath, serverVersion)
}

func getServerVersion(serverRootPath string) string {
	versionFile := serverRootPath + "/version.txt"
	fmt.Println(versionFile)
	exists, err := utils.FileExists(versionFile)
	if err != nil || !exists {
		return ""
	}
	data, err := os.ReadFile(versionFile)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func executeLuaScriptToExtractSettings(workDirPath, outputDirPath, serverVersion string) error {
	err := utils.MkDirIfNotExists(outputDirPath, 0755, true)
	if err != nil {
		return err
	}

	slog.Info("Extracting cluster json settings ...")
	sessionName := "dstm-extract-settings"
	cmd := "cd '" + workDirPath + "' && lua ./main.lua '" + outputDirPath + "' '" + serverVersion + "'"
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
	slog.Info("Cluster json settings extracted to " + outputDirPath)
	return nil
}

func prepareFiles(zipFilePath, tmpDirPath, scriptDirPath, destDirPath string) error {
	slog.Info("Unzipping scripts to " + tmpDirPath)
	err := utils.Unzip(zipFilePath, tmpDirPath, 0755)
	if err != nil {
		return err
	}

	slog.Info("Copying scripts ...")
	err = copyParserAndMocks(destDirPath)
	if err != nil {
		panic(err)
	}

	err = copyServerFiles(scriptDirPath, destDirPath)
	if err != nil {
		panic(err)
	}

	return nil
}

func copyParserAndMocks(destDirPath string) error {
	err := fs.WalkDir(scriptsDir, "worldgen_template", func(path string, d fs.DirEntry, err error) error {
		destPath := strings.Replace(path, "worldgen_template", destDirPath, 1)
		if d.IsDir() {
			return utils.MkDirIfNotExists(destPath, 0755, true)
		} else {
			file, err := scriptsDir.Open(path)
			if err != nil {
				return err
			}
			newFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			io.Copy(newFile, file)
		}
		return nil
	})
	return err
}

func copyServerFiles(scriptDirPath, destDirPath string) error {
	targetDirs := []string{
		scriptDirPath + "/languages",
		scriptDirPath + "/map/levels",
		scriptDirPath + "/map/tasksets",
	}
	targetFiles := []string{
		scriptDirPath + "/map/customize.lua",
		scriptDirPath + "/map/level.lua",
		scriptDirPath + "/map/levels.lua",
		scriptDirPath + "/map/locations.lua",
		scriptDirPath + "/map/resource_substitution.lua",
		scriptDirPath + "/map/settings.lua",
		scriptDirPath + "/map/startlocations.lua",
		scriptDirPath + "/map/tasksets.lua",

		scriptDirPath + "/constants.lua",
		scriptDirPath + "/strings.lua",

		scriptDirPath + "/class.lua",
		scriptDirPath + "/strict.lua",
		scriptDirPath + "/translator.lua",

		scriptDirPath + "/beefalo_clothing.lua",
		scriptDirPath + "/clothing.lua",
		scriptDirPath + "/emote_items.lua",
		scriptDirPath + "/item_blacklist.lua",
		scriptDirPath + "/misc_items.lua",
		scriptDirPath + "/prefabskins.lua",
		scriptDirPath + "/skin_strings.lua",
		scriptDirPath + "/techtree.lua",
		scriptDirPath + "/tuning.lua",
		scriptDirPath + "/worldsettings_overrides.lua",
	}
	targetPatterns := []string{
		scriptDirPath + "/speech_*.lua",
	}
	for _, pattern := range targetPatterns {
		pathList, err := filepath.Glob(pattern)
		if err != nil {
			return nil
		}
		targetFiles = append(targetFiles, pathList...)
	}

	for _, dirPath := range targetDirs {
		dest := strings.Replace(dirPath, scriptDirPath, destDirPath, 1)
		err := utils.CopyDir(dirPath, dest)
		if err != nil {
			panic(err)
		}
	}
	for _, filePath := range targetFiles {
		dest := strings.Replace(filePath, scriptDirPath, destDirPath, 1)
		err := utils.CopyFile(filePath, dest)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
