package web

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
	"time"
)

type Model map[string]interface{}
type M Model

var Funcs template.FuncMap = template.FuncMap{
	"getFormError": func(m map[string]string, key string) string {
		return m[key]
	},
}

type TmplCache struct {
	cache map[string]int64
	tmpls *template.Template
	sync.RWMutex
}

func NewTmplCache() *TmplCache {
	tc := &TmplCache{}
	tc.loadTemplates()
	return tc
}

func (tc *TmplCache) loadTemplates() {
	tc.Lock()
	tc.tmpls = nil
	tc.cache = make(map[string]int64, 0)
	tc.tmpls = template.Must(template.New("*").Funcs(Funcs).ParseGlob("templates/*.tmpl"))
	if f, _ := filepath.Glob("templates/*/*.tmpl"); len(f) > 0 {
		tc.tmpls = template.Must(tc.tmpls.ParseGlob("templates/*/*.tmpl"))
	}
	now := time.Now().Unix()
	for _, tmpl := range tc.tmpls.Templates() {
		tc.cache[tmpl.Name()] = now
	}
	tc.Unlock()
}

func (tc *TmplCache) Render(w http.ResponseWriter, r *http.Request, name string, vals map[string]interface{}) {
	tc.RLock()
	for name, modified := range tc.cache {
		modt := modtime(name)
		if modt == -1 {
			http.Redirect(w, r, "/error/404", 303)
			return
		}
		if modt > modified {
			tc.RUnlock()
			tc.loadTemplates()
			tc.RLock()
			break
		}
	}
	tc.RUnlock()
	if vals == nil {
		vals = make(map[string]interface{})
	}
	if msgk, msgv := GetFlash(w, r); msgk != "" {
		vals[msgk] = msgv
	}
	vals["session"] = GetAllSess(r)
	vals["formErrors"] = GetFormErrors(w, r)
	// modified to support template errors properly
	buf := new(bytes.Buffer)
	if err := tc.tmpls.ExecuteTemplate(buf, name, vals); err != nil {
		log.Printf("web >> tmpl.go >> Render() >> ExecuteTemplate() >> %v\n", err)
		http.Redirect(w, r, "/error/404", 303)
		return
	}
	buf.WriteTo(w)
	/*if err := tc.tmpls.ExecuteTemplate(w, name, vals); err != nil {
		http.Redirect(w, r, "/error/404", 303)
	}*/
	return
}

func (tc *TmplCache) RenderError(w http.ResponseWriter, r *http.Request, name string, code int) bool {
	tc.RLock()
	for name, modified := range tc.cache {
		modt := modtime(name)
		if modt > modified {
			tc.RUnlock()
			tc.loadTemplates()
			tc.RLock()
			break
		}
	}
	tc.RUnlock()
	vals := make(map[string]interface{})
	if msgk, msgv := GetFlash(w, r); msgk != "" {
		vals[msgk] = msgv
	}
	vals["session"] = GetAllSess(r)
	vals["formErrors"] = GetFormErrors(w, r)
	vals["code"] = code
	vals["error"] = http.StatusText(code)
	// modified to support template errors properly
	buf := new(bytes.Buffer)
	if err := tc.tmpls.ExecuteTemplate(buf, name, vals); err != nil {
		log.Printf("web >> tmpl.go >> Render() >> ExecuteTemplate() >> %v\n", err)
		return false
	}
	buf.WriteTo(w)
	return true
}

func modtime(name string) int64 {
	fp := "templates/" + name
	if f, _ := filepath.Glob("templates/*/" + name); len(f) == 1 {
		fp = f[0]
	}
	info, err := os.Stat(fp)
	if err != nil {
		return -1
	}
	return info.ModTime().Unix()
}
