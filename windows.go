//go:build windows
// +build windows

package platformdirs // import go.gideaworx.io/platformdirs

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

type fetchType uint8

const (
	FetchFromDLL fetchType = iota
	FetchFromRegistry
	FetchFromEnv
)

var (
	shell32          = syscall.NewLazyDLL("shell32.dll")
	shGetFolderPathW = shell32.NewProc("SHGetFolderPathW")

	varMap = map[string]map[fetchType]any{
		"CSIDL_LOCAL_APPDATA": {
			FetchFromDLL:      28,
			FetchFromRegistry: "Local AppData",
			FetchFromEnv:      "LOCALAPPDATA",
		},
	}
)

const (
	userDataDir = "CSIDL_LOCAL_APPDATA"
)

type windowsPlatformDirs struct {
	appAuthor string
	appName   string
	fetchType fetchType
}

func New(appAuthor, appName, _ string) PlatformDirs {
	return windowsPlatformDirs{
		appAuthor: appAuthor,
		appName:   appName,
		fetchType: getFetchType(),
	}
}

func getFetchType() fetchType {
	if shGetFolderPathW != nil {
		return FetchFromDLL
	}

	if _, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE); err == nil {
		return FetchFromRegistry
	}

	return FetchFromEnv
}

func (w windowsPlatformDirs) UserDataDir() (string, error) {
	return getDLLDir(userDataDir)
}

func (w windowsPlatformDirs) UserConfigDir() (string, error) {
	return w.UserDataDir()
}

func getDLLDir(varname string) (string, error) {
	var val = varMap[varname][FetchFromDLL].(int)
	var out uintptr

	_, _, err := shGetFolderPathW.Call(0, uintptr(val), 0, 0, out)

	noErr := syscall.Errno(0)
	if err != noErr {
		return "", err
	}

	return fmt.Sprint(out), nil
}

func getRegistryDir(_ string) (string, error) {
	return "", nil
}

func getEnvDir(_ string) (string, error) {
	return "", nil
}

func getDefaultDir(_ string) string {
	return filepath.Clean(os.Getenv("USERPROFILE"))
}
