package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/cagnosolutions/web"
)

var projectShareNew = web.Route{"POST", "/project/share/:name", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	r.ParseForm()
	var sharedProject SharedProject
	web.FormToStruct(&sharedProject, r.Form, "")
	project := r.FormValue(":name")
	if project == "" {
		ajaxResponse(w, `{"error":true,"output":"Error sharing project"}`)
		return
	}

	sharedProject.Created = time.Now().UnixNano()
	sharedProject.Id = genId()
	db.Set("sharedProject", sharedProject.Id, sharedProject)

	ajaxResponse(w, `{"error":false,"output":"Successfully started share","shareId":"`+sharedProject.Id+`"}`)
	return
}}

var projectShareView = web.Route{"GET", "/project/share/:name/:shareId", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	var sharedProject SharedProject
	if !db.Get("sharedProject", r.FormValue(":shareId"), &sharedProject) {
		//web.SetErrorRedirect(w, r, "/project", "You are not allowed to access this shared project")
		return
	}
	/*if !sharedProject.IsInvited(user.Username, user.Email) {

		web.SetErrorRedirect(w, r, "/project", "You are not allowed to access this shared project")
		return
	}*/
	tmpl.Render(w, r, "projectSharedView.tmpl", web.Model{
		"user":    user,
		"project": r.FormValue(":name"),
		"shareId": r.FormValue(":shareId"),
		"themes":  themes,
	})
	return
}}

var projectShareFiles = web.Route{"GET", "/project/share/:name/files/:shareId", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	var sharedProject SharedProject
	if !db.Get("sharedProject", r.FormValue(":shareId"), &sharedProject) {
		fmt.Println("Error getting project", r.FormValue(":shareId"))
		ajaxResponse(w, `{"id":"#","Children":false}`)
		return
	}
	/*if !sharedProject.IsInvited(user.Username, user.Email) {
		ajaxResponse(w, `{"id":"#","Children":false}`)
		return
	}*/
	path, _ := url.QueryUnescape(r.FormValue("path"))

	prePath := "./projects/" + sharedProject.UserId + "/" + r.FormValue(":name")

	files, err := ioutil.ReadDir(prePath + "/" + path)
	if err != nil {
		log.Printf("projectRoutes.go >> projectFiles >> ioutil.ReadDir() >> %v\n\n", err)
		ajaxResponse(w, `{"id":"#","children":false}`)
		return
	}
	fileNodes := []FileNode{}
	for _, file := range files {
		if file.Name()[0] != '.' {
			fileNode := FileNode{}
			filePath := path + "/" + file.Name()
			fileNode.Id = url.QueryEscape(filePath)
			fileNode.Text = file.Name()
			if file.IsDir() {
				fileNode.Type = "dir"
				fileNode.Children = !IsEmptyDir(prePath + "/" + filePath)
				fileNode.State = "closed"
			} else {
				fileNode.Type = setFileType(file.Name())
			}
			fileNodes = append(fileNodes, fileNode)
		}
	}
	b, err := json.Marshal(fileNodes)
	respString := string(b)
	if err != nil {
		log.Printf("projectRoutes.go >> projectFiles >> json.Marshal() >> %v\n\n", err)
		respString = `{"id":"#","Children":false}`
	}
	ajaxResponse(w, respString)
	return
}}

var projectShareFolderNew = web.Route{"POST", "/project/share/:name/mkdir/:shareId", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	var sharedProject SharedProject
	if !db.Get("sharedProject", r.FormValue(":shareId"), &sharedProject) {
		ajaxResponse(w, `{"error":true,"output":"You are not allowed to access this shared project"}`)
		return
	}
	/*if !sharedProject.IsInvited(user.Username, user.Email) {
		ajaxResponse(w, `{"error":true,"output":"You are not allowed to access this shared project"}`)
		return
	}*/
	if r.FormValue(":name") == "" {
		ajaxResponse(w, `{"error":true,"output":"Error creating new folder"}`)
		return
	}
	if r.FormValue("folder") == "" {

		ajaxResponse(w, `{"error":true,"output":"Error creating new folder"}`)
		return
	}
	p, _ := url.QueryUnescape(r.FormValue("path"))
	if p == "#" {
		p = ""
	}
	path := "projects/" + sharedProject.UserId + "/" + r.FormValue(":name") + "/" + p

	if err := os.MkdirAll(path+"/"+r.FormValue("folder"), 0755); err != nil {
		ajaxResponse(w, `{"error":true,"output":"Error creating new folder"}`)
		return
	}
	ajaxResponse(w, `{"error":false,"output":"Successfully created new folder"}`)
	return
}}

var projectShareFileNew = web.Route{"POST", "/project/share/:name/addFile/:shareId", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	var sharedProject SharedProject
	if !db.Get("sharedProject", r.FormValue(":shareId"), &sharedProject) {
		ajaxResponse(w, `{"error":true,"output":"You are not allowed to access this shared project"}`)
		return
	}
	/*if !sharedProject.IsInvited(user.Username, user.Email) {
		ajaxResponse(w, `{"error":true,"output":"You are not allowed to access this shared project"}`)
		return
	}*/
	if r.FormValue(":name") == "" {
		ajaxResponse(w, `{"error":true,"output":"Error creating new file"}`)
		return
	}
	if r.FormValue("file") == "" {
		ajaxResponse(w, `{"error":true,"output":"Error creating new file"}`)
		return
	}
	p, _ := url.QueryUnescape(r.FormValue("path"))
	if p == "#" {
		p = ""
	}
	path := "projects/" + sharedProject.UserId + "/" + r.FormValue(":name") + "/" + p

	f, err := os.Create(path + "/" + r.FormValue("file"))
	if err != nil {
		ajaxResponse(w, `{"error":true,"output":"Error creating new file"}`)
		return
	}
	f.Close()

	ajaxResponse(w, `{"error":false,"output":"Successfully created new file"}`)
	return
}}

var projectShareUploadImage = web.Route{"POST", "/project/share/:name/upload/:shareId", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	var sharedProject SharedProject
	if !db.Get("sharedProject", r.FormValue(":shareId"), &sharedProject) {
		ajaxResponse(w, `{"error":true,"output":"You are not allowed to access this shared project"}`)
		return
	}
	/*if !sharedProject.IsInvited(user.Username, user.Email) {
		ajaxResponse(w, `{"error":true,"output":"You are not allowed to access this shared project"}`)
		return
	}*/

	p, _ := url.QueryUnescape(r.FormValue("path"))
	if p == "#" {
		p = ""
	}
	path := "projects/" + sharedProject.UserId + "/" + r.FormValue(":name") + p

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("projectRoutes.go >> projectUploadImage >> r.FormFile() >> %v\n", err)
		ajaxResponse(w, `{"error":true,"output":"Error uploading image"}`)
		return
	}

	defer file.Close()
	f, err := os.OpenFile(path+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("projectRoutes.go >> projectUploadImage >> os.OpenFile() >> %v\n", err)
		ajaxResponse(w, `{"error":true,"output":"Error uploading image"}`)
		return
	}

	defer f.Close()
	if _, err := io.Copy(f, file); err != nil {
		log.Printf("projectRoutes.go >> projectUploadImage >> io.Copy() >> %v\n", err)
		ajaxResponse(w, `{"error":true,"output":"Error uploading image"}`)
		return
	}

	ajaxResponse(w, `{"error":false,"output":"Successfully uploaded image"}`)
	return
}}

var projectShareFileDel = web.Route{"POST", "/project/share/:name/file/del/:shareId", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	var sharedProject SharedProject
	if !db.Get("sharedProject", r.FormValue(":shareId"), &sharedProject) {
		ajaxResponse(w, `{"error":true,"output":"You are not allowed to access this shared project"}`)
		return
	}
	/*if !sharedProject.IsInvited(user.Username, user.Email) {
		ajaxResponse(w, `{"error":true,"output":"You are not allowed to access this shared project"}`)
		return
	}*/
	if r.FormValue("path") == "" {
		ajaxResponse(w, `{"error":true,"output":"Error deleting file/folder"}`)
		return
	}

	p, _ := url.QueryUnescape(r.FormValue("path"))
	path := "projects/" + sharedProject.UserId + "/" + r.FormValue(":name") + p

	if err := os.RemoveAll(path); err != nil {
		ajaxResponse(w, `{"error":true,"output":"Error deleting file/folder"}`)
		return
	}
	ajaxResponse(w, `{"error":false,"output":"Successfully deleted file/folder"}`)
	return
}}

var projectShareFileMove = web.Route{"POST", "/project/share/:name/file/move/:shareId", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	var sharedProject SharedProject
	if !db.Get("sharedProject", r.FormValue(":shareId"), &sharedProject) {
		ajaxResponse(w, `{"error":true,"output":"You are not allowed to access this shared project"}`)
		return
	}
	/*if !sharedProject.IsInvited(user.Username, user.Email) {
		ajaxResponse(w, `{"error":true,"output":"You are not allowed to access this shared project"}`)
		return
	}*/

	path := "projects/" + sharedProject.UserId + "/" + r.FormValue(":name")
	from, _ := url.QueryUnescape(r.FormValue("from"))
	to, _ := url.QueryUnescape(r.FormValue("to"))
	if from == "" || to == "" {
		log.Printf("projectRoutes.go >> projectFileMove >> FROM: %s, TO: %s\n\n", from, to)
		ajaxResponse(w, `{"error":true,"output":"Error `+r.FormValue("type")+`ing file/folder}`)
		return
	}

	if err := os.Rename(path+from, path+to); err != nil {
		log.Printf("projectRoutes.go >> projectFileMove >> os.Rename() >> %v\n\n", err)
		ajaxResponse(w, `{"error":true,"output":"Error `+r.FormValue("type")+`ing file/folder}`)
		return
	}
	ajaxResponse(w, `{"error":false,"output":"Successfully `+r.FormValue("type")+`ed file/folder"}`)
	return
}}

var projectShareFile = web.Route{"GET", "/project/share/:name/file/:shareId", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	var sharedProject SharedProject
	if !db.Get("sharedProject", r.FormValue(":shareId"), &sharedProject) {
		ajaxResponse(w, `{"error":true,"output":"You are not allowed to access this shared project"}`)
		return
	}
	/*if !sharedProject.IsInvited(user.Username, user.Email) {
		ajaxResponse(w, `{"error":true,"output":"You are not allowed to access this shared project"}`)
		return
	}*/
	p, _ := url.QueryUnescape(r.FormValue("path"))
	path := "projects/" + sharedProject.UserId + "/" + r.FormValue(":name") + p
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("projectRoutes.go >> projectFile >> ioutil.ReadFile() >> %v\n\n", err)
		ajaxResponse(w, `{"error":true,"output":"Error finding file"}`)
		return
	}

	ajaxResponse(w, fmt.Sprintf(`{"error":false,"output":"%s","fileType":"%s"}`, encodeFile(file), setFileType(path)))
	return
}}

var projectShareFileSave = web.Route{"POST", "/project/share/:name/file/save/:shareId", func(w http.ResponseWriter, r *http.Request) {

	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	var sharedProject SharedProject
	if !db.Get("sharedProject", r.FormValue(":shareId"), &sharedProject) {
		ajaxResponse(w, `{"error":true,"output":"You are not allowed to access this shared project"}`)
		return
	}
	/*if !sharedProject.IsInvited(user.Username, user.Email) {
		ajaxResponse(w, `{"error":true,"output":"You are not allowed to access this shared project"}`)
		return
	}*/

	p, _ := url.QueryUnescape(r.FormValue("path"))
	path := "projects/" + sharedProject.UserId + "/" + r.FormValue(":name") + p

	data := r.FormValue("data")

	if err := ioutil.WriteFile(path, []byte(data), 0666); err != nil {
		log.Printf("projectRoutes.go >> projectFileSave >> ioutil.WriteFile() >> %v\n\n", err)
		ajaxResponse(w, `{"error":true,"output":"Error saving `+filepath.Base(path)+`"}`)
		return
	}

	ajaxResponse(w, `{"error":false,"output":"Successfully saved `+filepath.Base(path)+`"}`)
	return
}}
