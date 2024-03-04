package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func ExpandPath(path string) string {
	p := os.ExpandEnv(path)
	if strings.HasPrefix(p, "~") {
		p = strings.Replace(p, "~", os.Getenv("HOME"), 1)
	}
	return p
}

// Build a absolute path string that is enclosed in double quotes
func BuildAbsPathString(comps ...string) (string, error) {
	if !filepath.IsAbs(comps[0]) {
		return "", errors.New("path must be absolute")
	}
	path := filepath.Join(comps...)
	path = strings.ReplaceAll(path, "\"", "\\\"")
	return "\"" + path + "\"", nil
}
