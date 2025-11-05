package main

import (
	"os"
)

func checkPath(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	return err
}

func isFile(path string) bool {
	info, _ := os.Stat(path)
	return info.Mode().IsRegular()
}
