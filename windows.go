//go:build windows
// +build windows

package platformdirs

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

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
	shell32,           = syscall.NewLazyDLL("shell32.dll")
	shGetFolderPathW   = shell32.NewProc("SHGetFolderPathW")

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

func (w windowsPlatformDirs) UserDataDir() (string, error) {
	return getDLLDir(userDataDir)
}

func (w windowsPlatformDirs) UserConfigDir() (string, error) {
	return w.UserDataDir()
}

func getDLLDir(varname string) (string, error) {
	var val = varMap[varname][FetchFromDLL].(int)
	var out uintptr

	r1, _, err := shGetFolderPathW.Call(0, uintptr(val), 0, 0, out)

	noErr := syscall.Errno(0)
	if err != noErr {
		return "", err
	}

	return fmt.Sprint(out), nil
}

func getRegistryDir(varname string) (string, error) {
	return "", nil
}

func getEnvDir(varname string) (string, error) {
	return "", nil
}

func getDefaultDir(varname string) string {
	return filepath.Clean(os.Getenv("USERPROFILE"))
}

func uintptrToString(u uintptr) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Print("panicked in uintptrToString")
			return
		}
	}()

	ptr := (*uint16)(unsafe.Pointer(u))
	return windows.UTF16PtrToString(ptr)
}
