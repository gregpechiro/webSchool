package main

type User struct {
	Id        string `json:"id"`
	Email     string `json:"email,omitempty" auth:"username" required:"register,login,account"`
	Password  string `json:"password,omitempty" auth:"password" required:"register,login"`
	Active    bool   `json:"active" auth:"active"`
	Role      string `json:"role,omitempty"`
	FirstName string `json:"firstName,omitempty" required:"register,account"`
	LastName  string `json:"lastName,omitempty" required:"register,account"`
	Created   int64  `json:"created,omitempty"`
	LastSeen  int64  `json:"lastSeen,omitempty"`
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
