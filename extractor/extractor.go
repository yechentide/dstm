package extractor

import (
	"embed"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yechentide/dstm/shell"
	"github.com/yechentide/dstm/utils"
)

//go:embed scripts
var scriptsDir embed.FS

func ExtractSettings(zipFile, tmpDir string) error {
	zipFilePath := utils.ExpandPath(zipFile)
	tmpDirPath := utils.ExpandPath(tmpDir)

	err := prepareFiles(zipFilePath, tmpDirPath)
	if err != nil {
		return err
	}

	workDir := tmpDirPath + "/work"
	outputDir := tmpDirPath + "/output"
	return executeLuaScriptToExtractSettings(workDir, outputDir)
}

func executeLuaScriptToExtractSettings(workDir, outputDir string) error {
	err := utils.MkDirIfNotExists(outputDir, 0755, true)
	if err != nil {
		return err
	}

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
	return nil
}

func prepareFiles(zipFile, tmpDir string) error {
	err := utils.Unzip(zipFile, tmpDir, 0755)
	if err != nil {
		return err
	}

	workDir := tmpDir + "/work"
	err = utils.MkDirIfNotExists(workDir, 0755, true)
	if err != nil {
		panic(err)
	}

	err = copyParserAndMocks(tmpDir)
	if err != nil {
		panic(err)
	}

	err = copyServerFiles(tmpDir)
	if err != nil {
		panic(err)
	}

	return nil
}

func copyParserAndMocks(tmpDir string) error {
	err := fs.WalkDir(scriptsDir, "scripts", func(path string, d fs.DirEntry, err error) error {
		if d.Name() == "scripts" {
			return nil
		}
		destPath := strings.Replace(path, "scripts", tmpDir+"/work", 1)
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

func copyServerFiles(tmpDir string) error {
	targetDirs := []string{
		tmpDir + "/scripts/languages",
		tmpDir + "/scripts/map/levels",
		tmpDir + "/scripts/map/tasksets",
	}
	targetFiles := []string{
		tmpDir + "/scripts/map/customize.lua",
		tmpDir + "/scripts/map/level.lua",
		tmpDir + "/scripts/map/levels.lua",
		tmpDir + "/scripts/map/locations.lua",
		tmpDir + "/scripts/map/resource_substitution.lua",
		tmpDir + "/scripts/map/settings.lua",
		tmpDir + "/scripts/map/startlocations.lua",
		tmpDir + "/scripts/map/tasksets.lua",

		tmpDir + "/scripts/constants.lua",
		tmpDir + "/scripts/strings.lua",

		tmpDir + "/scripts/class.lua",
		tmpDir + "/scripts/strict.lua",
		tmpDir + "/scripts/translator.lua",

		tmpDir + "/scripts/beefalo_clothing.lua",
		tmpDir + "/scripts/clothing.lua",
		tmpDir + "/scripts/emote_items.lua",
		tmpDir + "/scripts/item_blacklist.lua",
		tmpDir + "/scripts/misc_items.lua",
		tmpDir + "/scripts/prefabskins.lua",
		tmpDir + "/scripts/skin_strings.lua",
		tmpDir + "/scripts/techtree.lua",
		tmpDir + "/scripts/tuning.lua",
		tmpDir + "/scripts/worldsettings_overrides.lua",
	}
	targetPatterns := []string{
		tmpDir + "/scripts/speech_*.lua",
	}
	for _, pattern := range targetPatterns {
		pathList, err := filepath.Glob(pattern)
		if err != nil {
			return nil
		}
		targetFiles = append(targetFiles, pathList...)
	}

	for _, dirPath := range targetDirs {
		dest := strings.Replace(dirPath, tmpDir+"/scripts", tmpDir+"/work", 1)
		err := utils.CopyDir(dirPath, dest)
		if err != nil {
			panic(err)
		}
	}
	for _, filePath := range targetFiles {
		dest := strings.Replace(filePath, tmpDir+"/scripts", tmpDir+"/work", 1)
		err := utils.CopyFile(filePath, dest)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
