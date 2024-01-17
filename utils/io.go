package utils

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
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

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	p := ExpandPath(filepath)
	out, err := os.Create(p)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func FileExists(path string) (bool, error) {
	p := ExpandPath(path)
	f, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	if f.IsDir() {
		return false, errors.New("file is a directory: " + p)
	}
	return true, nil
}

func DirExists(path string) (bool, error) {
	p := ExpandPath(path)
	f, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	if !f.IsDir() {
		return false, errors.New("file is not a directory: " + p)
	}
	return true, nil
}

func MkDirIfNotExists(path string, perm fs.FileMode, recursive bool) error {
	p := ExpandPath(path)
	exists, err := DirExists(p)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	if recursive {
		return os.MkdirAll(p, perm)
	}
	return os.Mkdir(p, perm)
}

func DelDirIfExists(path string) error {
	p := ExpandPath(path)
	exists, err := DirExists(p)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	return os.RemoveAll(p)
}
