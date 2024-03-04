package utils

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// https://qiita.com/brushwood-field/items/417f7c07ee5813239ff3

func Unzip(srcPath, destPath string, mode fs.FileMode) error {
	r, err := zip.OpenReader(srcPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Mode().IsDir() {
			if err := os.MkdirAll(filepath.Join(destPath, f.Name), mode); err != nil {
				return err
			}
			continue
		}
		if err := saveUnzippedFile(destPath, *f); err != nil {
			return err
		}
	}
	return nil
}

func saveUnzippedFile(destDirPath string, f zip.File) error {
	destPath := filepath.Join(destDirPath, f.Name)

	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, rc); err != nil {
		return err
	}
	return nil
}
