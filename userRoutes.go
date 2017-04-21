package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var temp = web.Route{"GET", "/temp", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "temp.tmpl", nil)
}}

/* --- All Project management --- */

var project = web.Route{"GET", "/project", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	model := web.Model{
		"user": user,
	}
	var projects []FileStats
	path := "projects/" + id

	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		model["alertError"] = "Error getting projects"
	} else {
		for _, file := range files {
			if file.IsDir() {
				size, numFiles, _ := DirStats(path + "/" + file.Name())
				projects = append(projects, FileStats{
					Name:     file.Name(),
					Size:     PrettySize(size),
					NumFiles: numFiles,
				})

			}
		}

	}
	model["projects"] = projects
	tmpl.Render(w, r, "project.tmpl", model)
	return
}}

var projectNew = web.Route{"POST", "/project", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	name := r.FormValue("name")
	if name == "" {
		web.SetErrorRedirect(w, r, "/project", "Project name cannot be empty")
		return
	}
	path := "projects/" + id + "/" + name
	_, err := os.Stat(path)
	if err == nil {
		web.SetErrorRedirect(w, r, "/project", "That project name is already taken")
		return
	}
	if os.IsExist(err) {
		web.SetErrorRedirect(w, r, "/project", "Error creating new project")
		return
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		web.SetErrorRedirect(w, r, "/project", "Error creating new project")
		return
	}

	web.SetSuccessRedirect(w, r, "/project", "Successfully created project")
	return
}}

/* --- Account management --- */

var account = web.Route{"GET", "/account", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	tmpl.Render(w, r, "account.tmpl", web.Model{
		"user": user,
	})
	return
}}

var accountSave = web.Route{"POST", "/account", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}

	r.ParseForm()
	if r.FormValue("password") == "" {
		r.Form.Set("password", user.Password)
	}
	if errs, ok := web.FormToStruct(&user, r.Form, "account"); !ok {
		web.SetFormErrors(w, errs)
		web.SetErrorRedirect(w, r, "/account", "Error updating account information")
		return
	}

	// check for uniqueness
	var users []User
	db.TestQuery("user", &users, adb.Eq("email", user.Email), adb.Ne("id", `"`+user.Id+`"`))
	if len(users) > 0 {
		web.SetErrorRedirect(w, r, "/account", "Error updating account information.\nEmail is already in use.")
		return
	}

	db.Set("user", user.Id, user)
	web.SetSuccessRedirect(w, r, "/account", "Successfully updated information")

}}
