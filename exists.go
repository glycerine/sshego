package sshego

import (
	"os"
)

// fileExists returns true iff the path name is a file (and not a directory or non-existant).
func fileExists(name string) bool {
	fi, err := os.Stat(name)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return false
	}
	return true
}

// FileExistsLen check if name is an actually file (directories don't count) and also
// returns the length of the file.
func fileExistsLen(name string) (bool, int64) {
	fi, err := os.Stat(name)
	if err != nil {
		return false, 0
	}
	if fi.IsDir() {
		return false, 0
	}
	return true, fi.Size()
}

// DirExists returns true if name represents a directory on disk.
func dirExists(name string) bool {
	fi, err := os.Stat(name)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return true
	}
	return false
}
