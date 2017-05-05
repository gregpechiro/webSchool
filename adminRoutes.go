package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var adminHome = web.Route{"GET", "/admin", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
		return
	}
	tmpl.Render(w, r, "admin.tmpl", web.Model{
		"user": user,
	})
	return
}}

var adminUser = web.Route{"GET", "/admin/user", func(w http.ResponseWriter, r *http.Request) {
	var users []User
	db.All("user", &users)
	tmpl.Render(w, r, "user.tmpl", web.Model{
		"users": users,
		"user":  User{},
	})
	return
}}

var adminUserOne = web.Route{"GET", "/admin/user/:id", func(w http.ResponseWriter, r *http.Request) {
	var users []User
	db.All("user", &users)

	var user User
	if !db.Get("user", r.FormValue(":id"), &user) {
		web.SetErrorRedirect(w, r, "/admin/user", "Error finding user")
		return
	}

	tmpl.Render(w, r, "user.tmpl", web.Model{
		"users": users,
		"user":  user,
	})
	return
}}

var adminUserSave = web.Route{"POST", "/admin/user", func(w http.ResponseWriter, r *http.Request) {
	var user User
	id := r.FormValue("id")
	db.Get("user", id, &user)

	r.ParseForm()
	if r.FormValue("password") == "" {
		if user.Id == "" {
			web.SetFormErrors(w, map[string]string{
				"adminUser.password": "Password is required",
			})
			web.SetErrorRedirect(w, r, "/admin/user", "Error saving user")
			return
		}
		r.Form.Set("password", user.Password)
	}

	redirect := "/admin/user"
	if id != "" {
		redirect += "/" + id
	}

	if errs, ok := web.FormToStruct(&user, r.Form, "adminUser"); !ok {
		web.SetFormErrors(w, errs)
		web.SetErrorRedirect(w, r, redirect, "Error updating account information")
		return
	}

	// check for uniqueness
	var users []User
	db.TestQuery("user", &users, adb.Eq("email", user.Email), adb.Ne("id", `"`+user.Id+`"`))
	if len(users) > 0 {
		web.SetErrorRedirect(w, r, redirect, "Error saving user information.<br>Email is already in use.")
		return
	}

	db.TestQuery("user", &users, adb.Eq("username", user.Username), adb.Ne("id", `"`+user.Id+`"`))
	if len(users) > 0 {
		web.SetErrorRedirect(w, r, redirect, "Error saving user information.<br>Username is already in use.")
		return
	}

	if user.Id == "" {
		user.Id = genId()
		user.Created = time.Now().Unix()
		// create user project folder
		if err := os.MkdirAll("projects/"+user.Id, 0755); err != nil {
			web.SetErrorRedirect(w, r, "/admin/user", "Error creating user")
			return
		}
	}

	db.Set("user", user.Id, user)
	web.SetSuccessRedirect(w, r, "/admin/user/"+user.Id, "Successfully saved user")
	return
}}

var adminUserDel = web.Route{"POST", "/admin/user/:id", func(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue(":id")
	if err := os.RemoveAll("projects/" + id); err != nil {
		log.Printf("adminRoutes.go >> adminUserDel >> os.RemoveAll() >> %v\n", err)
		web.SetErrorRedirect(w, r, "/admin/user/"+id, "Error deleting user")
		return
	}

	db.Del("user", id)
	web.SetSuccessRedirect(w, r, "/admin/user", "Successfully deleted user")
	return
}}
