package env

import (
	"errors"
	"log/slog"
	"runtime"

	"github.com/yechentide/dstm/logger"
	"gopkg.in/ini.v1"
)

type OSHelper interface {
	IsRoot() bool
	HasSudoPermmission() bool
	Dependencies() []string
	IsInstalled([]string) (map[string]bool, error)
	IsTerminalMultiplexerReady() (bool, error)
	InstallPackages([]string, string) error
	InstallAllRequired(string) error
}

func checkOS() {
	current := runtime.GOOS
	if runtime.GOOS != "linux" {
		slog.Error("Unsupported OS: " + current)
		logger.PrintJsonResultAndExit(1)
	}

	osInfo, err := ini.Load("/etc/os-release")
	if err != nil {
		slog.Error("Failed to load /etc/os-release", "error", err)
		logger.PrintJsonResultAndExit(1)
	}
	distroID := osInfo.Section("").Key("ID").String()
	value, found := supportedOS[distroID]
	if found {
		slog.Debug("Detected distro: " + distroID)
	} else {
		slog.Error("Unsupported distro: " + distroID)
		logger.PrintJsonResultAndExit(1)
	}
	osDistro = value
	// osVer = osInfo.Section("").Key("VERSION_ID").String()
}

func GetOSHelper() (OSHelper, error) {
	checkOS()
	switch osDistro {
	case debian:
		return &DebianHelper{}, nil
	default:
		return nil, errors.New("unsupported OS")
	}
}
