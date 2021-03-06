package main

type User struct {
	Id        string `json:"id"`
	Email     string `json:"email,omitempty" auth:"username" required:"register,login,account,adminUser"`
	Password  string `json:"password,omitempty" auth:"password" required:"register,login"`
	Active    bool   `json:"active" auth:"active"`
	Role      string `json:"role,omitempty"`
	FirstName string `json:"firstName,omitempty" required:"register,account,adminUser"`
	LastName  string `json:"lastName,omitempty" required:"register,account,adminUser"`
	Username  string `json:"username" required:"register,account,adminUser"`
	Created   int64  `json:"created,omitempty"`
	LastSeen  int64  `json:"lastSeen,omitempty"`
}

type SharedProject struct {
	Id          string   `json:"id"`
	UserId      string   `json:"userId"`
	ProjectName string   `json:"projectName"`
	Invites     []string `json:"invites"`
	URL         string   `json:"url"`
	Created     int64    `json:"created"`
}

func (s SharedProject) IsInvited(username, email string) bool {
	for _, n := range s.Invites {
		if n == username || n == email {
			return true
		}
	}
	return false
}

type FileStats struct {
	Name     string
	Size     string
	NumFiles int
}

type FileNode struct {
	Id       string `json:"id"`
	Text     string `json:"text"`
	Children bool   `json:"children"`
	Type     string `json:"type"`
	State    string `json:"state"`
}
