package utils

import (
	"embed"
	"errors"
	"io"
	"os"
)

func FileExists(filePath string) (bool, error) {
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
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.WriteString(destFile, content)
	return err
}

func CopyFile(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

func CopyEmbeddedFile(embeddedDir embed.FS, filePath, destPath string) error {
	srcFile, err := embeddedDir.Open(filePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}
