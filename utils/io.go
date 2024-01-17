package utils

import (
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func ExpandPath(path string) string {
	p := os.ExpandEnv(path)
	if strings.HasPrefix(p, "~/") {
		p = strings.Replace(p, "~/", os.Getenv("HOME")+"/", 1)
	}
	return p
}

func FileExists(path string) (bool, error) {
	p := ExpandPath(path)
	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func MkDir(path string, perm fs.FileMode, recursive bool) error {
	exists, err := FileExists(path)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	p := ExpandPath(path)
	if recursive {
		return os.MkdirAll(p, 0755)
	}
	return os.Mkdir(p, 0755)
}
