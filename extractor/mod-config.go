package extractor

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/yechentide/dstm/global"
	"github.com/yechentide/dstm/shell"
	"github.com/yechentide/dstm/utils"
)

func ExtractModConfiguration(modDirPath, outputDirPath, langCode string) error {
	isModDir, err := utils.IsModDir(modDirPath)
	if err != nil {
		return err
	}
	if !isModDir {
		return errors.New(modDirPath + " is not a mod directory")
	}

	modID, err := utils.GetModIDFromPath(modDirPath)
	if err != nil {
		return err
	}

	luaScriptsDirPath, err := prepareLuaScripts()
	if err != nil {
		return err
	}

	err = os.MkdirAll(outputDirPath, 0755)
	if err != nil {
		return err
	}

	slog.Info("Extracting mod configuration ...")

	sessionName := global.SESSION_EXTRACT_MOD_CONFIG
	cmd := "cd '" + luaScriptsDirPath + "' && lua extract-mod-config.lua '" + modDirPath + "' '" + outputDirPath + "' '" + langCode + "'"
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

	jsonPath := outputDirPath + "/" + langCode + "." + modID + ".json"
	exists, err := utils.FileExists(jsonPath)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("failed to extract mod configuration")
	}
	slog.Info("Mod configuration extracted to " + jsonPath)
	return nil
}
