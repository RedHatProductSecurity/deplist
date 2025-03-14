package utils

import (
	"archive/zip"
	"path/filepath"
)

func CharIsDigit(c string) bool {
	if len(c) == 0 {
		return false
	}

	if c[0] < '0' || c[0] > '9' {
		return false
	}
	return true
}

// Quickly scan for java files within the zip file
func ZipContainsJava(path string) (bool, error) {
	read, err := zip.OpenReader(path)
	if err != nil {
		return false, err
	}
	defer read.Close()

	for _, f := range read.File {
		switch filepath.Ext(f.Name) {
		case ".java":
			fallthrough
		case ".war":
			fallthrough
		case ".ear":
			fallthrough
		case ".jar":
			return true, nil
		}
	}

	return false, nil
}
