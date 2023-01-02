package domain

type User struct {
	ID       int64
	Username string
	Password string
	Name     string
	Locale   string
	Blocked  bool
	Admin    bool
}
