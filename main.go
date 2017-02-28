package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var tmpl *web.TmplCache
var mux *web.Mux
var db *adb.DB = adb.NewDB()

func init() {
	db.AddStore("user")

	web.SESSDUR = 15 * time.Minute
	mux = web.NewMux()

	// unsecure routes
	mux.AddRoutes(home, register, login, logout, loginPost)

	// user routes
	mux.AddSecureRoutes(USER, account)

	// admin routes
	mux.AddSecureRoutes(ADMIN, adminHome)

	tmpl = web.NewTmplCache()

	defaultUsers()
}

func main() {
	fmt.Println(">>> DID YOU REGISTER ANY NEW ROUTES <<<")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

var home = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "home.tmpl", nil)
	return
}}
