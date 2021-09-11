package domain

type Version string

func (v Version) String() string {
	return string(v)
}

type Secret string

func (s Secret) String() string {
	return string(s)
}
