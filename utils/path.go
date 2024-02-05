package utils

import (
	"os"
	"strings"
)

func ExpandPath(path string) string {
	p := os.ExpandEnv(path)
	if strings.HasPrefix(p, "~/") {
		p = strings.Replace(p, "~/", os.Getenv("HOME")+"/", 1)
	}
	return p
}
