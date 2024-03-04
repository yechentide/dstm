package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

func GetModIDFromPath(modDirPath string) (string, error) {
	pattern := `^(workshop-)?(\d+)$`
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	dirName := filepath.Base(modDirPath)
	matched := reg.FindStringSubmatch(dirName)
	if len(matched) != 3 {
		return "", fmt.Errorf("invalid mod directory name: %s", dirName)
	}

	return matched[len(matched)-1], nil
}
