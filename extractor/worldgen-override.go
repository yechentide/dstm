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

func ExtractWorldgenOverride(shardDirPath, outputDirPath string) error {
	luaScriptsDirPath, err := prepareLuaScripts()
	if err != nil {
		return err
	}

	err = os.MkdirAll(outputDirPath, 0755)
	if err != nil {
		return err
	}

	slog.Info("Extracting worldgen override ...")

	sessionName := global.SESSION_EXTRACT_WORLD_OVERRIDE
	cmd := "cd '" + luaScriptsDirPath + "' && lua convert-worldgen-override.lua '" + shardDirPath + "' '" + outputDirPath + "'"
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

	jsonPath := outputDirPath + "/worldgenoverride.json"
	exists, err := utils.FileExists(jsonPath)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("failed to extract worldgen override")
	}
	slog.Info("Worldgen override extracted to " + outputDirPath)
	return nil
}
