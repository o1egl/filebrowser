//go:generate ../../tools/bin/go-enum --sql --marshal --nocase --names --file $GOFILE

package domain

/*
ENUM(
fs
)
*/
type VolumeType int

type Volume struct {
	ID          int64
	Type        VolumeType
	Label       string
	Path        string
	Description string
}

type VolumePermissions struct {
	Read   bool
	Create bool
	Modify bool
	Delete bool
	Share  bool
}

type VolumeWithPermissions struct {
	Volume      Volume
	Permissions VolumePermissions
}
