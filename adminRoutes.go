package main

import (
	"net/http"

	"github.com/cagnosolutions/web"
)

var adminHome = web.Route{"GET", "/admin", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error retrieving user")
		return
	}
	tmpl.Render(w, r, "admin.tmpl", web.Model{
		"user": user,
	})
}}
