package utils

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func DirExists(dirPath string) (bool, error) {
	f, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	if !f.IsDir() {
		return false, errors.New("file is not a directory: " + dirPath)
	}
	return true, nil
}

func MkDirIfNotExists(dirPath string, perm fs.FileMode, recursive bool) error {
	exists, err := DirExists(dirPath)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	if recursive {
		return os.MkdirAll(dirPath, perm)
	}
	return os.Mkdir(dirPath, perm)
}

func DelDirIfExists(dirPath string) error {
	exists, err := DirExists(dirPath)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	return os.RemoveAll(dirPath)
}

func RemakeDir(dirPath string, perm fs.FileMode, recursive bool) error {
	err := DelDirIfExists(dirPath)
	if err != nil {
		return err
	}
	return MkDirIfNotExists(dirPath, perm, recursive)
}

func CopyDir(srcPath, destPath string) error {
	srcDir, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	err = MkDirIfNotExists(destPath, srcDir.Mode(), true)
	if err != nil {
		return err
	}

	err = filepath.Walk(srcPath, func(path string, info fs.FileInfo, err error) error {
		if path == srcPath {
			return nil
		}
		var e error
		destPath := strings.Replace(path, srcPath, destPath, 1)
		if info.IsDir() {
			e = MkDirIfNotExists(destPath, info.Mode(), true)
		} else {
			e = CopyFile(path, destPath)
		}
		if e != nil {
			return e
		}
		return nil
	})
	return err
}

// list directory names
func ListDirs(parentDirPath string) ([]string, error) {
	files, err := os.ReadDir(parentDirPath)
	if err != nil {
		return nil, nil
	}
	var dirs []string
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}
	return dirs, nil
}
