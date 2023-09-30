package platformdirs

type PlatformDirs interface {
	UserDataDir() (string, error)
	UserConfigDir() (string, error)
}
