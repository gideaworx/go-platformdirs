//go:build darwin
// +build darwin

package platformdirs

import (
	"os"
	"path/filepath"
)

type darwinPlatformDirs struct {
	appName    string
	appVersion string
}

func (l darwinPlatformDirs) UserDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library", "Application Support", l.appName, l.appVersion), nil
}

func (l darwinPlatformDirs) UserConfigDir() (string, error) {
	return l.UserDataDir()
}

func New(_, appName, appVersion string) PlatformDirs {
	return darwinPlatformDirs{
		appName, appVersion,
	}
}
