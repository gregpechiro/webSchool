package main

type User struct {
	Id        string `json:"id"`
	Email     string `json:"email,omitempty" auth:"username" required:"register, login"`
	Password  string `json:"password,omitempty" auth:"password" required:"register, login"`
	Active    bool   `json:"active" auth:"active"`
	Role      string `json:"role,omitempty"`
	FirstName string `json:"firstName,omitempty" required:"register"`
	LastName  string `json:"lastName,omitempty" required:"register"`
	Created   int64  `json:"created,omitempty"`
	LastSeen  int64  `json:"lastSeen,omitempty"`
}
