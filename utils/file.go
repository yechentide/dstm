package utils

import (
	"errors"
	"io"
	"os"
)

func FileExists(path string) (bool, error) {
	filePath := ExpandPath(path)
	f, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	if f.IsDir() {
		return false, errors.New("file is a directory: " + filePath)
	}
	return true, nil
}

func WriteToFile(content, destPath string) error {
	destFile, err := os.Create(ExpandPath(destPath))
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.WriteString(destFile, content)
	return err
}

func CopyFile(srcPath, destPath string) error {
	srcFile, err := os.Open(ExpandPath(srcPath))
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(ExpandPath(destPath))
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}
