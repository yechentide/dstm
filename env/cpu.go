package env

import (
	"log/slog"

	"github.com/yechentide/dstm/logger"
	"github.com/yechentide/dstm/shell"
)

func getCPUArch() cpuArch {
	bits, err := shell.ExecuteAndGetOutput("getconf", "LONG_BIT")
	if err != nil {
		slog.Error("Failed to get CPU arch", "error", err)
		logger.PrintJsonResultAndExit(1)
	}
	if bits == "32" {
		return bit32
	}
	return bit64
}

func Is64BitCPU() bool {
	return getCPUArch() == bit64
}
