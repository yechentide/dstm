package extractor

import (
	"embed"
	"errors"
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

//go:embed scripts
var scriptsDir embed.FS

func ExtractSettings(zipFile, outputDir string) error {
	zipFilePath := utils.ExpandPath(zipFile)
	tmpDir := "/tmp/dstm-extract-json"
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

	outputDir = utils.ExpandPath(outputDir)
	return executeLuaScriptToExtractSettings(workDir, outputDir)
}

func executeLuaScriptToExtractSettings(workDir, outputDir string) error {
	err := utils.MkDirIfNotExists(outputDir, 0755, true)
	if err != nil {
		return err
	}

	slog.Info("Extracting cluster json settings ...")
	sessionName := "dstm-extract-settings"
	cmd := "cd '" + workDir + "' && lua ./main.lua '" + workDir + "/languages' '" + outputDir + "'"
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
	exists, err := utils.FileExists(outputDir + "/en.forest.master.json")
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("failed to extract settings")
	}
	slog.Info("Cluster json settings extracted to " + outputDir)
	return nil
}

func prepareFiles(zipFile, tmpDir, scriptDir, destDir string) error {
	slog.Info("Unzipping scripts to " + tmpDir)
	err := utils.Unzip(zipFile, tmpDir, 0755)
	if err != nil {
		return err
	}

	slog.Info("Copying scripts ...")
	err = copyParserAndMocks(destDir)
	if err != nil {
		panic(err)
	}

	err = copyServerFiles(scriptDir, destDir)
	if err != nil {
		panic(err)
	}

	return nil
}

func copyParserAndMocks(destDir string) error {
	err := fs.WalkDir(scriptsDir, "scripts", func(path string, d fs.DirEntry, err error) error {
		destPath := strings.Replace(path, "scripts", destDir, 1)
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

func copyServerFiles(scriptDir, destDir string) error {
	targetDirs := []string{
		scriptDir + "/languages",
		scriptDir + "/map/levels",
		scriptDir + "/map/tasksets",
	}
	targetFiles := []string{
		scriptDir + "/map/customize.lua",
		scriptDir + "/map/level.lua",
		scriptDir + "/map/levels.lua",
		scriptDir + "/map/locations.lua",
		scriptDir + "/map/resource_substitution.lua",
		scriptDir + "/map/settings.lua",
		scriptDir + "/map/startlocations.lua",
		scriptDir + "/map/tasksets.lua",

		scriptDir + "/constants.lua",
		scriptDir + "/strings.lua",

		scriptDir + "/class.lua",
		scriptDir + "/strict.lua",
		scriptDir + "/translator.lua",

		scriptDir + "/beefalo_clothing.lua",
		scriptDir + "/clothing.lua",
		scriptDir + "/emote_items.lua",
		scriptDir + "/item_blacklist.lua",
		scriptDir + "/misc_items.lua",
		scriptDir + "/prefabskins.lua",
		scriptDir + "/skin_strings.lua",
		scriptDir + "/techtree.lua",
		scriptDir + "/tuning.lua",
		scriptDir + "/worldsettings_overrides.lua",
	}
	targetPatterns := []string{
		scriptDir + "/speech_*.lua",
	}
	for _, pattern := range targetPatterns {
		pathList, err := filepath.Glob(pattern)
		if err != nil {
			return nil
		}
		targetFiles = append(targetFiles, pathList...)
	}

	for _, dirPath := range targetDirs {
		dest := strings.Replace(dirPath, scriptDir, destDir, 1)
		err := utils.CopyDir(dirPath, dest)
		if err != nil {
			panic(err)
		}
	}
	for _, filePath := range targetFiles {
		dest := strings.Replace(filePath, scriptDir, destDir, 1)
		err := utils.CopyFile(filePath, dest)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
