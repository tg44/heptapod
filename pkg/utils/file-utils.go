package utils

import (
	"os/user"
	"path/filepath"
)

//https://stackoverflow.com/a/43578461/2118749
func FixupPathsToHandleHome(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}
