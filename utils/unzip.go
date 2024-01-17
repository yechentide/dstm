package utils

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// https://qiita.com/brushwood-field/items/417f7c07ee5813239ff3

func Unzip(src, dest string, perm fs.FileMode) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Mode().IsDir() {
			continue
		}
		if err := saveUnzippedFile(dest, *f, perm); err != nil {
			return err
		}
	}
	return nil
}

func saveUnzippedFile(destDir string, f zip.File, perm fs.FileMode) error {
	destPath := filepath.Join(destDir, f.Name)
	if err := os.MkdirAll(filepath.Dir(destPath), perm); err != nil {
		return err
	}

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
