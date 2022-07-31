package domain

type User struct {
	ID           int64
	Username     string
	Password     string
	Home         HomeVolume
	Name         string
	Locale       string
	LockPassword bool
	Blocked      bool
	IsAdmin      bool
}

type HomeVolume struct {
	Path        string
	Permissions VolumePermissions
}
