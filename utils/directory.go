package utils

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

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

func RemakeDir(path string, perm fs.FileMode, recursive bool) error {
	err := DelDirIfExists(path)
	if err != nil {
		return err
	}
	return MkDirIfNotExists(path, perm, recursive)
}

func CopyDir(src, dest string) error {
	srcPath := ExpandPath(src)
	destRootPath := ExpandPath(dest)

	srcDir, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	err = MkDirIfNotExists(destRootPath, srcDir.Mode(), true)
	if err != nil {
		return err
	}

	err = filepath.Walk(srcPath, func(path string, info fs.FileInfo, err error) error {
		if path == src {
			return nil
		}
		var e error
		destPath := strings.Replace(path, srcPath, destRootPath, 1)
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

func ListDirs(parentDir string) ([]string, error) {
	files, err := os.ReadDir(parentDir)
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
