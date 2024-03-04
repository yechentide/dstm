package utils

import (
	"embed"
	"errors"
	"io"
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

func RemakeDir(dirPath string, mode fs.FileMode, recursive bool) error {
	err := os.RemoveAll(dirPath)
	if err != nil {
		return err
	}
	if recursive {
		return os.MkdirAll(dirPath, mode)
	} else {
		return os.Mkdir(dirPath, mode)
	}
}

func CopyDir(srcPath, destPath string) error {
	srcDir, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(destPath, srcDir.Mode())
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
			e = os.MkdirAll(destPath, info.Mode())
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

// List directory names
func ListChildDirs(parentDirPath string) ([]string, error) {
	files, err := os.ReadDir(parentDirPath)
	if err != nil {
		return nil, err
	}
	var dirs []string
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}
	return dirs, nil
}

func CopyEmbeddedDir(embeddedDir embed.FS, root string, destDirPath string) error {
	err := fs.WalkDir(embeddedDir, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		destPath := strings.Replace(path, root, destDirPath, 1)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		file, err := embeddedDir.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		newFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		return err
	})
	return err
}
