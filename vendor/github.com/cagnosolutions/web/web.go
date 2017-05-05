package web

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	STATIC_PATH string = "static"

	ACTIVITY_URL = "/activity/manager/ftw"

	DEFAULT_ERR_ROUTE Route = Route{"GET", "/error/:code", func(w http.ResponseWriter, r *http.Request) {
		code, err := strconv.Atoi(r.FormValue(":code"))
		if err != nil {
			code = 500
		}

		vals := map[string]interface{}{
			"session": GetAllSess(r),
			"code":    code,
			"status":  http.StatusText(int(code)),
		}

		if t, err := template.ParseFiles(fmt.Sprintf("errors/error%d.tmpl", code)); err == nil {
			buf := new(bytes.Buffer)
			if err := t.Execute(buf, vals); err == nil {
				buf.WriteTo(w)
				return
			}
		}

		if t, err := template.ParseFiles("errors/error.tmpl"); err == nil {
			buf := new(bytes.Buffer)
			if err := t.Execute(buf, vals); err == nil {
				buf.WriteTo(w)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html; utf-8")
		fmt.Fprintf(w, "<html><body><head><title>%d</title></head><center><br/><h1>HTTP Status %d %s</h1><p>Default Error Handler</p></center></body></html>", code, code, http.StatusText(int(code)))
		return
	}}

	ACTIVITY_MANAGER_ROUTE Route = Route{"GET", ACTIVITY_URL, handleActivity}

	DEFAULT_HANDLER func(w http.ResponseWriter, r *http.Request) bool = func(w http.ResponseWriter, r *http.Request) bool {
		return false
	}
)

type Mux struct {
	handlers []*Handler
	static   http.Handler
}

func NewMux() *Mux {
	if _, err := os.Stat(STATIC_PATH); os.IsNotExist(err) {
		if err := os.MkdirAll(STATIC_PATH, 0755); err != nil {
			log.Fatalf("could not create static file path %q: %v\n", STATIC_PATH, err)
		}
	}
	mux := &Mux{
		handlers: make([]*Handler, 0),
		static:   http.StripPrefix("/static/", http.FileServer(http.Dir(STATIC_PATH))),
	}
	mux.Add(DEFAULT_ERR_ROUTE)
	if AMANAGER {
		mux.Add(ACTIVITY_MANAGER_ROUTE)
	}
	return mux
}

func (mux *Mux) GetRoutes() []*Route {
	var routes []*Route
	for _, handler := range mux.handlers {
		routes = append(routes, handler.route)
	}
	return routes
}

func (mux *Mux) Add(route Route) {
	mux.handlers = append(mux.handlers, &Handler{&route, nil})
}

func (mux *Mux) AddRoutes(routes ...Route) {
	for _, route := range routes {
		mux.Add(route)
	}
}

func (mux *Mux) AddSecure(auth Auth, route Route) {
	mux.handlers = append(mux.handlers, &Handler{&route, &auth})
}

func (mux *Mux) AddSecureRoutes(auth Auth, routes ...Route) {
	for _, route := range routes {
		mux.AddSecure(auth, route)
	}
}

func (mux *Mux) Get(path string, fn http.HandlerFunc) {
	mux.Add(Route{"GET", path, fn})
}

func (mux *Mux) GetSecure(auth Auth, path string, fn http.HandlerFunc) {
	mux.AddSecure(auth, Route{"GET", path, fn})
}

func (mux *Mux) Post(path string, fn http.HandlerFunc) {
	mux.Add(Route{"POST", path, fn})
}

func (mux *Mux) PostSecure(auth Auth, path string, fn http.HandlerFunc) {
	mux.AddSecure(auth, Route{"POST", path, fn})
}

func (mux *Mux) Put(path string, fn http.HandlerFunc) {
	mux.Add(Route{"PUT", path, fn})
}

func (mux *Mux) PutSecure(auth Auth, path string, fn http.HandlerFunc) {
	mux.AddSecure(auth, Route{"PUT", path, fn})
}

func (mux *Mux) Delete(path string, fn http.HandlerFunc) {
	mux.Add(Route{"DELETE", path, fn})
}

func (mux *Mux) DeleteSecure(auth Auth, path string, fn http.HandlerFunc) {
	mux.AddSecure(auth, Route{"DELETE", path, fn})
}

func (mux *Mux) Err(fn http.HandlerFunc) {
	mux.handlers[0].route.Fn = fn
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check default for completion
	if DEFAULT_HANDLER(w, r) {
		return
	}

	// ignore OPTIONS and favicon calls
	if r.Method == "OPTIONS" || r.URL.Path == "/favicon.ico" {
		return
	}
	// handle static content
	if r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/static/") {
		mux.static.ServeHTTP(w, r)
		return
	}
	// check for post request in order to validate via the referer
	if r.Method == "POST" && !strings.Contains(r.Referer(), r.Host) { // possibly add origin check in there too...
		// invalid request redirect to 403
		http.Redirect(w, r, "/error/403", 303)
		return
	}
	// otherwise, attempt to handle
	for _, handler := range mux.handlers {
		// if this handler.route;s method does not match...
		if handler.route.Method != r.Method {
			continue // ...continue/skip to next handler.route
		}
		// parse path vars and check for match (simultaneously)
		if vals, ok := parse(handler.route.Path, r.URL.Path); ok {
			if len(vals) > 0 {
				r.URL.RawQuery += ("&" + url.Values(vals).Encode())
			}
			if AMANAGER && handler.route.Path != ACTIVITY_URL && !strings.HasPrefix(r.URL.Path, "/error/") {
				actId := GetActivityId(r)
				if act, ok := ActivityManager.Get(actId); ok {
					act.Update(r)
					ActivityManager.Put(actId, act)
				} else {
					PutActivityId(w, ActivityManager.New(r))
				}
			}

			// we have a path match
			//check for security constraints...
			if handler.auth != nil {
				if Authorized(w, r, handler.auth.Roles) {
					handler.route.Fn(w, r)
					return
				}
				// we have a match, but the role is not currently allowed... so we redirect
				SetErrorRedirect(w, r, handler.auth.Redirect, handler.auth.Msg)
				return
			}
			// we have a match and our vals were aded to r (get with r.FormValue(":key"))
			handler.route.Fn(w, r)
			return
		}
	}
	// no matches found, return an error page
	http.Redirect(w, r, "/error/404", 303)
	return
}

type Handler struct {
	route *Route
	auth  *Auth
}

type Auth struct {
	Roles         []string
	Redirect, Msg string
}

type Route struct {
	Method string
	Path   string
	Fn     http.HandlerFunc
}

// parse registered pattern
func parse(routepath, path string) (url.Values, bool) {
	p := make(url.Values)
	var i, j int
	for i < len(path) {
		switch {
		case j >= len(routepath):
			if routepath != "/" && len(routepath) > 0 && routepath[len(routepath)-1] == '/' {
				return p, true
			}
			return nil, false
		case routepath[j] == ':':
			var name, val string
			var nextc byte
			name, nextc, j = match(routepath, isBoth, j+1)
			val, _, i = match(path, byteParse(nextc), i)
			p.Add(":"+name, val)
		case path[i] == routepath[j]:
			i++
			j++
		default:
			return nil, false
		}
	}
	if j != len(routepath) {
		return nil, false
	}
	return p, true
}

// match path with registered handler
func match(s string, f func(byte) bool, i int) (matched string, next byte, j int) {
	j = i
	for j < len(s) && f(s[j]) {
		j++
	}
	if j < len(s) {
		next = s[j]
	}
	return s[i:j], next, j
}

// determine type of byte
func byteParse(b byte) func(byte) bool {
	return func(c byte) bool {
		return c != b && c != '/'
	}
}

// test for alpha byte
func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// test for numerical byte
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// test for alpha or numerical byte
func isBoth(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}
