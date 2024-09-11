package platformdirs // import go.gideaworx.io/platformdirs

type PlatformDirs interface {
	UserDataDir() (string, error)
	UserConfigDir() (string, error)
}
