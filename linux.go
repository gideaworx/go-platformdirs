//go:build linux
// +build linux

package platformdirs

import (
	"os"
	"path/filepath"
	"strings"
)

type linuxPlatformDirs struct {
	appName    string
	appVersion string
}

func (l linuxPlatformDirs) UserDataDir() (string, error) {
	xdgHome := os.Getenv("XDG_DATA_HOME")
	if strings.TrimSpace(xdgHome) == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		xdgHome = filepath.Join(home, ".local", "share")
	}

	return filepath.Join(xdgHome, l.appName, l.appVersion), nil
}

func (l linuxPlatformDirs) UserConfigDir() (string, error) {
	return l.UserDataDir()
}

func New(_, appName, appVersion string) PlatformDirs {
	return linuxPlatformDirs{
		appName, appVersion,
	}
}
