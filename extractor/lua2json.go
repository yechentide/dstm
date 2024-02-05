package extractor

import (
	_ "embed"
	"errors"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/yechentide/dstm/shell"
	"github.com/yechentide/dstm/utils"
)

var (
	//go:embed lua2json/main.lua
	mainScript string
	//go:embed worldgen_template/my_utils.lua
	utilsScript string
)

func GenerateWorldgenOverrideJson(shardDirPath string) (string, error) {
	fileNameForServer := "worldgenoverride.lua"
	fileNameForPC := "leveldataoverride.lua"

	exists, err := utils.FileExists(shardDirPath + "/" + fileNameForServer)
	if err != nil {
		return "", err
	}
	if exists {
		return ConvertLuaObjectToJson(shardDirPath + "/" + fileNameForServer)
	}

	exists, err = utils.FileExists(shardDirPath + "/" + fileNameForPC)
	if err != nil {
		return "", err
	}
	if exists {
		return ConvertLuaObjectToJson(shardDirPath + "/" + fileNameForPC)
	}

	return "", errors.New("file not found: " + shardDirPath + "/" + fileNameForServer)
}

func ConvertLuaObjectToJson(luaFilePath string) (string, error) {
	tmpDir := "/tmp/dstm-lua-to-json"
	err := utils.RemakeDir(tmpDir, 0755, true)
	if err != nil {
		return "", err
	}

	err = utils.WriteToFile(mainScript, tmpDir+"/main.lua")
	if err != nil {
		return "", err
	}
	err = utils.WriteToFile(utilsScript, tmpDir+"/my_utils.lua")
	if err != nil {
		return "", err
	}

	fileNameWithExtension := filepath.Base(luaFilePath)
	fileName := strings.TrimSuffix(fileNameWithExtension, filepath.Ext(fileNameWithExtension))
	parentDirPath := filepath.Dir(luaFilePath)

	slog.Info("Converting lua object to json ...")
	sessionName := "dstm-lua-to-json"
	cmd := "cd '" + tmpDir + "' && lua ./main.lua '" + parentDirPath + "' '" + fileName + "'"
	err = shell.CreateTmuxSession(sessionName, cmd)
	if err != nil {
		return "", err
	}
	for {
		time.Sleep(1 * time.Second)
		sessionExists, err := shell.HasTmuxSession(sessionName)
		if err != nil {
			return "", err
		}
		if !sessionExists {
			break
		}
	}

	outputFilePath := parentDirPath + "/" + fileName + ".json"
	exists, err := utils.FileExists(outputFilePath)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errors.New("failed to extract settings")
	}
	slog.Info("Extracted json object to " + outputFilePath)
	return outputFilePath, nil
}
