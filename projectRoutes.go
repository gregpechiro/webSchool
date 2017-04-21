package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/cagnosolutions/web"
)

var projectView = web.Route{"GET", "/project/:name", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	tmpl.Render(w, r, "projectView.tmpl", web.Model{
		"user":    user,
		"project": r.FormValue(":name"),
		"themes":  themes,
	})
	return
}}

var projectRename = web.Route{"POST", "/project/:name", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	name := r.FormValue(":name")
	if name == "" {
		web.SetErrorRedirect(w, r, "/project", "Project name cannot be empty")
		return
	}
	newName := r.FormValue("newName")
	if name == "" {
		web.SetErrorRedirect(w, r, "/project", "Project name cannot be empty")
		return
	}

	if name == newName {
		web.SetErrorRedirect(w, r, "/project", "New project name must be different than original project name")
	}

	path := "projects/" + id + "/" + name
	newPath := "projects/" + id + "/" + newName

	if _, err := os.Stat(path); err != nil {
		web.SetErrorRedirect(w, r, "/project", "Error finding project")
		return
	}

	_, err := os.Stat(newPath)
	if err == nil {
		web.SetErrorRedirect(w, r, "/project", "That project name is already taken")
		return
	}
	if os.IsNotExist(err) {
		web.SetErrorRedirect(w, r, "/project", "Error creating new project")
		return
	}

	if err := os.Rename(path, newPath); err != nil {
		web.SetErrorRedirect(w, r, "/project", "Error renaming project")
		return
	}

	web.SetSuccessRedirect(w, r, "/project", "Successfully renamed project")
	return
}}

var projectFiles = web.Route{"GET", "/project/:name/files", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	path, _ := url.QueryUnescape(r.FormValue("path"))

	prePath := "./projects/" + id + "/" + r.FormValue(":name")

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

var projectFolderNew = web.Route{"POST", "/project/:name/mkdir", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	if r.FormValue(":name") == "" {
		// web.SetErrorRedirect(w, r, "/project/"+r.FormValue(":name"), "Error creating new folder")
		ajaxResponse(w, `{"error":true,"output":"Error creating new folder"}`)
		return
	}
	if r.FormValue("folder") == "" {

		// web.SetErrorRedirect(w, r, "/project/"+r.FormValue(":name"), "Error creating new folder")
		ajaxResponse(w, `{"error":true,"output":"Error creating new folder"}`)
		return
	}
	p, _ := url.QueryUnescape(r.FormValue("path"))
	if p == "#" {
		p = ""
	}
	path := "projects/" + id + "/" + r.FormValue(":name") + "/" + p

	if err := os.MkdirAll(path+"/"+r.FormValue("folder"), 0755); err != nil {
		// web.SetErrorRedirect(w, r, "/project/"+r.FormValue(":name"), "Error creating new folder")
		ajaxResponse(w, `{"error":true,"output":"Error creating new folder"}`)
		return
	}
	// web.SetSuccessRedirect(w, r, "/project/"+r.FormValue(":name"), "Successfully created new folder")
	ajaxResponse(w, `{"error":false,"output":"Successfully created new folder"}`)
	return
}}

var projectFileNew = web.Route{"POST", "/project/:name/addFile", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	if r.FormValue(":name") == "" {
		// web.SetErrorRedirect(w, r, "/project/"+r.FormValue(":name"), "Error creating file file")
		ajaxResponse(w, `{"error":true,"output":"Error creating new file"}`)
		return
	}
	if r.FormValue("file") == "" {
		// web.SetErrorRedirect(w, r, "/project/"+r.FormValue(":name"), "Error creating new file")
		ajaxResponse(w, `{"error":true,"output":"Error creating new file"}`)
		return
	}
	p, _ := url.QueryUnescape(r.FormValue("path"))
	if p == "#" {
		p = ""
	}
	path := "projects/" + id + "/" + r.FormValue(":name") + "/" + p

	f, err := os.Create(path + "/" + r.FormValue("file"))
	if err != nil {
		// web.SetErrorRedirect(w, r, "/project/"+r.FormValue(":name"), "Error creating new file")
		ajaxResponse(w, `{"error":true,"output":"Error creating new file"}`)
		return
	}
	f.Close()

	// web.SetSuccessRedirect(w, r, "/project/"+r.FormValue(":name"), "Successfully created new file")
	ajaxResponse(w, `{"error":false,"output":"Successfully created new file"}`)
	return
}}

var projectFileDel = web.Route{"POST", "/project/:name/file/del", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	if r.FormValue("path") == "" {
		ajaxResponse(w, `{"error":true,"output":"Error deleting file/folder"}`)
		return
	}

	p, _ := url.QueryUnescape(r.FormValue("path"))
	path := "projects/" + id + "/" + r.FormValue(":name") + p

	if err := os.RemoveAll(path); err != nil {
		ajaxResponse(w, `{"error":true,"output":"Error deleting file/folder"}`)
		return
	}
	ajaxResponse(w, `{"error":false,"output":"Successfully deleted file/folder"}`)
	return
}}

var projectFileMove = web.Route{"POST", "/project/:name/file/move", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}

	path := "projects/" + id + "/" + r.FormValue(":name")
	from, _ := url.QueryUnescape(r.FormValue("from"))
	to, _ := url.QueryUnescape(r.FormValue("to"))
	if from == "" || to == "" {
		log.Printf("projectRoutes.go >> projectFileMove >> FROM: %s, TO: %s\n\n", from, to)
		// web.SetErrorRedirect(w, r, "/project/"+r.FormValue(":name"), "Error "+r.FormValue("type")+"ing file/folder")
		ajaxResponse(w, `{"error":true,"output":"Error `+r.FormValue("type")+`ing file/folder}`)
		return
	}

	if err := os.Rename(path+from, path+to); err != nil {
		log.Printf("projectRoutes.go >> projectFileMove >> os.Rename() >> %v\n\n", err)
		// web.SetErrorRedirect(w, r, "/project/"+r.FormValue(":name"), "Error "+r.FormValue("type")+"ing file/folder")
		ajaxResponse(w, `{"error":true,"output":"Error `+r.FormValue("type")+`ing file/folder}`)
		return
	}
	// web.SetSuccessRedirect(w, r, "/project/"+r.FormValue(":name"), "Successfully "+r.FormValue("type")+"ed file/folder")
	ajaxResponse(w, `{"error":false,"output":"Successfully `+r.FormValue("type")+`ed file/folder"}`)
	return
}}

var projectFile = web.Route{"GET", "/project/:name/file", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}
	p, _ := url.QueryUnescape(r.FormValue("path"))
	path := "projects/" + id + "/" + r.FormValue(":name") + p
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("projectRoutes.go >> projectFile >> ioutil.ReadFile() >> %v\n\n", err)
		ajaxResponse(w, `{"error":true,"output":"Error finding file"}`)
		return
	}

	ajaxResponse(w, fmt.Sprintf(`{"error":false,"output":"%s","fileType":"%s"}`, encodeFile(file), setFileType(path)))
	return
}}

var projectFileSave = web.Route{"POST", "/project/:name/file/save", func(w http.ResponseWriter, r *http.Request) {
	id := web.GetId(r)
	var user User
	if !db.Get("user", id, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding user")
	}

	p, _ := url.QueryUnescape(r.FormValue("path"))
	path := "projects/" + id + "/" + r.FormValue(":name") + p

	data := r.FormValue("data")

	if err := ioutil.WriteFile(path, []byte(data), 0666); err != nil {
		log.Printf("projectRoutes.go >> projectFileSave >> ioutil.WriteFile() >> %v\n\n", err)
		ajaxResponse(w, `{"error":true,"output":"Error saving `+filepath.Base(path)+`"}`)
		return
	}

	ajaxResponse(w, `{"error":false,"output":"Successfully saved `+filepath.Base(path)+`"}`)
	return
}}
