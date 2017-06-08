package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var tmpl *web.TmplCache
var mux *web.Mux
var db *adb.DB = adb.NewDB()

const HOST = "http://localhost:9999"

// const HOST = "http://school.xiphoid24.com"

func init() {
	db.AddStore("user")
	db.AddStore("sharedProject")

	web.SESSDUR = 45 * time.Minute

	web.DEFAULT_HANDLER = func(w http.ResponseWriter, r *http.Request) bool {
		subs := strings.Split(r.Host, ".")
		if len(subs) == 4 {
			username := subs[0]
			var user User
			if db.TestQueryOne("user", &user, adb.Eq("username", username)) {
				http.FileServer(http.Dir("projects/"+user.Id)).ServeHTTP(w, r)
				return true
			}
			http.Redirect(w, r, HOST+"/error/404", 303)
			return true
		}
		return false
	}

	mux = web.NewMux()

	// unsecure routes
	mux.AddRoutes(home, register, login, logout, loginPost, updateSession, temp)

	// user routes
	mux.AddSecureRoutes(USER, project, projectNew, projectDel, account, accountSave)

	mux.AddSecureRoutes(USER, projectRename, projectView, projectFiles, projectFile, projectFolderNew)
	mux.AddSecureRoutes(USER, projectFileNew, projectFileDel, projectFileMove, projectFileSave, projectUploadImage)

	mux.AddSecureRoutes(USER, projectShareNew, projectShareView, projectShareFiles, projectShareFile, projectShareFolderNew)
	mux.AddSecureRoutes(USER, projectShareFileNew, projectShareFileDel, projectShareFileMove, projectShareFileSave, projectShareUploadImage)

	// admin routes
	mux.AddSecureRoutes(ADMIN, adminHome, adminUser, adminUserOne, adminUserSave, adminUserDel)

	web.Funcs["pretty"] = pretty
	web.Funcs["prettyDateTime"] = PrettyDateTime
	tmpl = web.NewTmplCache()

	defaultUsers()

}

func main() {
	fmt.Println(">>> DID YOU REGISTER ANY NEW ROUTES <<<")
	log.Fatal(http.ListenAndServe(":9999", mux))
}

var home = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "home.tmpl", nil)
	return
}}
