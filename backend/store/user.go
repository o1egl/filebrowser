package store

// User holds user-related info
type User struct {
	ID       string `json:"id" storm:"id"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Provider string `json:"provider"`
	// username and password are used only by local provider
	Username     string      `json:"username" storm:"index"`
	Password     string      `json:"password"`
	Scope        string      `json:"scope"`
	Locale       string      `json:"locale"`
	Rules        []Rule      `json:"rules"`
	Commands     []string    `json:"commands"`
	LockPassword bool        `json:"lockPassword"`
	Permissions  Permissions `json:"permissions"`
	Blocked      bool        `json:"blocked"`
}

// Rule is a allow/disallow rule.
type Rule struct {
	Regex bool   `json:"regex"`
	Allow bool   `json:"allow"`
	Rule  string `json:"rule"`
}

// Permissions describe a user's permissions.
type Permissions struct {
	Admin    bool `json:"admin"`
	Execute  bool `json:"execute"`
	Create   bool `json:"create"`
	Rename   bool `json:"rename"`
	Modify   bool `json:"modify"`
	Delete   bool `json:"delete"`
	Share    bool `json:"share"`
	Download bool `json:"download"`
}

func (p Permissions) CanExecute() bool {
	return p.Execute || p.Admin
}

func (p Permissions) CanCreate() bool {
	return p.Create || p.Admin
}

func (p Permissions) CanRename() bool {
	return p.Rename || p.Admin
}

func (p Permissions) CanModify() bool {
	return p.Modify || p.Admin
}

func (p Permissions) CanDelete() bool {
	return p.Delete || p.Admin
}

func (p Permissions) CanShare() bool {
	return p.Share || p.Admin
}

func (p Permissions) CanDownload() bool {
	return p.Download || p.Admin
}
