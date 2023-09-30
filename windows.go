//go:build windows
// +build windows

package platformdirs

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

type fetchType uint8

const (
	FetchFromDLL fetchType = iota
	FetchFromRegistry
	FetchFromEnv
)

var (
	shell32, _          = syscall.LoadLibrary("shell32.dll")
	shGetFolderPathW, _ = syscall.GetModuleAddress(shell32, "SHGetFolderPathW")

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
	if shGetFolderPathW != 0 {
		return FetchFromDLL
	}

	if _, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE); err == nil {
		return FetchFromRegistry
	}

	return FetchFromEnv
}

func (w windowsPlatformDirs) getUserDataDir() (string, error) {
	return getDLLDir(userDataDir)
}

func getDLLDir(varname string) (string, error) {
	var val = varMap[varname][FetchFromDLL].(int)
	var out uintptr

	r1, r2, errptr := syscall.SyscallN(getShortPathNameW, 0, val, 0, 0, out)

	fmt.Println("test")
	fmt.Println(windows.UTF16PtrToString(r1))
	fmt.Println(windows.UTF16PtrToString(r2))
	fmt.Println(windows.UTF16PtrToString(out))
	return "", nil
}

func getRegistryDir(varname string) (string, error) {

}

func getEnvDir(varname string) (string, error) {

}

func getDefaultDir(varname sring) string {
	return filepath.Clean(os.Getenv("USERPROFILE"))
}
