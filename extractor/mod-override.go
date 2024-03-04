package extractor

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/yechentide/dstm/shell"
	"github.com/yechentide/dstm/utils"
)

func ExtractModOverride(shardDirPath, outputDirPath string) error {
	luaScriptsDirPath, err := prepareLuaScripts()
	if err != nil {
		return err
	}

	err = os.MkdirAll(outputDirPath, 0755)
	if err != nil {
		return err
	}

	slog.Info("Extracting mod override ...")

	sessionName := "dstm-extract-mod-override"
	cmd := "cd '" + luaScriptsDirPath + "' && lua convert-mod-override.lua '" + shardDirPath + "' '" + outputDirPath + "'"
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

	jsonPath := outputDirPath + "/modoverrides.json"
	exists, err := utils.FileExists(jsonPath)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("failed to extract mod override")
	}
	slog.Info("Mod override extracted to " + jsonPath)
	return nil
}
