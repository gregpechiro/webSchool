package web

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/* ~ = ~ = ~ = ~ = ~ = *==[ COOKIE LIBRARY ]==* ~ = ~ = ~ = ~ = ~ = */
var EXPYEAR = 3
var EXPMONTH = 0
var EXPDAY = 0

func Expires() time.Time {
	return time.Now().AddDate(EXPYEAR, EXPMONTH, EXPDAY)
}

type Cook struct {
	name, value string
}

func GetCookie(r *http.Request, name string) string {
	cookie, err := getCookie(r, name)
	if err != nil || cookie == nil {
		return ""
	}
	return BaseDec(cookie.Value)
}

func PutCookie(w http.ResponseWriter, name, value string) {
	cookie := FreshCookie(name, value, Expires())
	http.SetCookie(w, &cookie)
}

func DeleteCookie(w http.ResponseWriter, name string) {
	cookie := FreshCookie(name, "", time.Now())
	cookie.MaxAge = -1
	http.SetCookie(w, &cookie)
}

func getCookie(r *http.Request, name string) (*http.Cookie, error) {
	return r.Cookie(UrlEnc(name))
}

func putCookie(w http.ResponseWriter, c *http.Cookie) {
	http.SetCookie(w, c)
}

func deleteCookie(w http.ResponseWriter, c *http.Cookie) {
	c.MaxAge = -1
	http.SetCookie(w, c)
}

func GetAll(r *http.Request) []Cook {
	return getAll(r.Cookies())
}

func getAll(cookies []*http.Cookie) []Cook {
	var cooks []Cook
	for _, c := range cookies {
		cooks = append(cooks, Cook{UrlDec(c.Name), Dec(c.Value)})
	}
	return cooks
}

func GetStartsWith(r *http.Request, pre string) []Cook {
	var cookies []*http.Cookie
	for _, cookie := range r.Cookies() {
		if strings.HasPrefix(UrlDec(cookie.Name), pre) {
			cookies = append(cookies, cookie)
		}
	}
	return getAll(cookies)
}

func UpdateStartsWith(w http.ResponseWriter, r *http.Request, pre string, exp time.Time) {
	for _, cookie := range r.Cookies() {
		if strings.HasPrefix(UrlDec(cookie.Name), pre) {
			cookie.Expires = exp
			cookie.Path = "/"
			cookie.HttpOnly = true
			http.SetCookie(w, cookie)
		}
	}
}

func DeleteStartsWith(w http.ResponseWriter, r *http.Request, pre string) {
	for _, cookie := range r.Cookies() {
		if strings.HasPrefix(UrlDec(cookie.Name), pre) {
			cookie.Expires = time.Now()
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
		}
	}
}

func GetActivityId(r *http.Request) string {
	return GetCookie(r, "ACTID")
}

func PutActivityId(w http.ResponseWriter, aid string) {
	PutCookie(w, "ACTID", aid)
}

func FreshCookie(name, value string, expires time.Time) http.Cookie {
	return http.Cookie{
		Name:     UrlEnc(name),
		Value:    BaseEnc(value),
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
	}
}

/* ~ = ~ = ~ = ~ = ~ = *==[ FLASH VARIABES UTILIZING COOKIE LIB]==* ~ = ~ = ~ = ~ = ~ = */
func SetFlash(w http.ResponseWriter, kind, msg string) {
	PutCookie(w, "flash", kind+":"+msg)
}

func GetFlash(w http.ResponseWriter, r *http.Request) (string, string) {
	msg := strings.Split(GetCookie(r, "flash"), ":")
	DeleteCookie(w, "flash")
	if len(msg) != 2 {
		return "", ""
	}
	return msg[0], msg[1]
}

func SetFlashRedirect(w http.ResponseWriter, r *http.Request, url, kind, msg string) {
	SetFlash(w, kind, msg)
	http.Redirect(w, r, url, 303)
	return
}

func SetSuccessRedirect(w http.ResponseWriter, r *http.Request, url, msg string) {
	SetFlash(w, "alertSuccess", msg)
	http.Redirect(w, r, url, 303)
	return
}

func SetErrorRedirect(w http.ResponseWriter, r *http.Request, url, msg string) {
	SetFlash(w, "alertError", msg)
	http.Redirect(w, r, url, 303)
	return
}

func SetMsgRedirect(w http.ResponseWriter, r *http.Request, url, msg string) {
	SetFlash(w, "alert", msg)
	http.Redirect(w, r, url, 303)
	return
}

func SetFormErrors(w http.ResponseWriter, errors map[string]string) {
	b, err := json.Marshal(errors)
	if err != nil {
		log.Printf("Web >> cook.go >> SetFormErrors(): %v \n", err)
		return
	}
	PutCookie(w, "formErrors", string(b))
}

func GetFormErrors(w http.ResponseWriter, r *http.Request) map[string]string {
	var m map[string]string
	data := GetCookie(r, "formErrors")
	if len(data) > 0 && data[0] == '{' {
		if err := json.Unmarshal([]byte(data), &m); err != nil {
			log.Printf("Web >> cook.go >> GetFormErrors(): %v \n", err)
		}
	}
	DeleteCookie(w, "formErrors")
	return m
}

/* ~ = ~ = ~ = ~ = ~ = *==[ SESSION MGMT UTILIZING COOKIE LIB ]==* ~ = ~ = ~ = ~ = ~ = */
type session map[string]interface{}

func (s session) PutId(w http.ResponseWriter, id string) {
	s["ID"] = id
	cookie := FreshCookie("SESS-D", toString(s), SessDur())
	putCookie(w, &cookie)
	return
}

func (s session) GetRole() string {
	role := s["ROLE"]
	if role == nil {
		return ""
	}
	if r, ok := role.(string); ok {
		return r
	}
	return ""
}

func (s session) GetId() string {
	id := s["ID"]
	if id == nil {
		return ""
	}
	if i, ok := id.(string); ok {
		return i
	}
	return ""
}

var SESSDUR = time.Minute * 15
var COOKIESALT = 50

func NewCookieSalt() {
	rand.Seed(time.Now().UnixNano())
	COOKIESALT = rand.Intn(100)
}

func toMap(data string) session {
	var m map[string]interface{}
	data = Dec(data)
	if len(data) > 0 && data[0] == '{' {
		if err := json.Unmarshal([]byte(data), &m); err != nil {
			log.Printf("Web >> cook.go >> toMap(): %v \n", err)
		}
	}
	return m
}

func toString(m map[string]interface{}) string {
	b, err := json.Marshal(m)
	if err != nil {
		log.Printf("Web >> cook.go >> toString(): %v \n", err)
		return ""
	}
	return Enc(string(b))
}

func SessDur() time.Time {
	return time.Now().Add(SESSDUR)
}

func Login(w http.ResponseWriter, r *http.Request, role string) session {
	m := toMap(GetCookie(r, "SESS-D"))
	if m == nil {
		m = make(map[string]interface{})
	}
	if role == "" {
		return m
	}
	if _, ok := m["ROLE"]; ok {
		return m
	}
	m["ROLE"] = role
	cookie := FreshCookie("SESS-D", toString(m), SessDur())
	putCookie(w, &cookie)
	return m
}

func Logout(w http.ResponseWriter) {
	DeleteCookie(w, "SESS-D")
}

func GetRole(r *http.Request) string {
	role := GetSess(r, "ROLE")
	if role == nil {
		return ""
	}
	if r, ok := role.(string); ok {
		return r
	}
	return ""
}

func GetId(r *http.Request) string {
	id := GetSess(r, "ID")
	if id == nil {
		return ""
	}
	if i, ok := id.(string); ok {
		return i
	}
	return ""
}

func PutSess(w http.ResponseWriter, r *http.Request, name string, val interface{}) {
	m := toMap(GetCookie(r, "SESS-D"))
	if m == nil {
		m = make(map[string]interface{})
	}
	m[name] = val
	cookie := FreshCookie("SESS-D", toString(m), SessDur())
	putCookie(w, &cookie)
}

func GetSess(r *http.Request, key string) interface{} {
	m := toMap(GetCookie(r, "SESS-D"))
	return m[key]
}

func GetAllSess(r *http.Request) session {
	return toMap(GetCookie(r, "SESS-D"))
}

func PutMultiSess(w http.ResponseWriter, r *http.Request, n map[string]interface{}) {
	m := toMap(GetCookie(r, "SESS-D"))
	var cookie http.Cookie
	if m == nil {
		m = make(map[string]interface{})
	}
	if len(m) < len(n) {
		for k, v := range m {
			n[k] = v
		}
		cookie = FreshCookie("SESS-D", toString(n), SessDur())
	} else {
		for k, v := range n {
			m[k] = v
		}
		cookie = FreshCookie("SESS-D", toString(m), SessDur())
	}
	putCookie(w, &cookie)
}

func Authorized(w http.ResponseWriter, r *http.Request, reqRoles []string) bool {
	m := toMap(GetCookie(r, "SESS-D"))
	role := m["ROLE"]
	if m == nil || role == "" {
		DeleteCookie(w, "SESS-D")
		return false
	}
	cookie := FreshCookie("SESS-D", toString(m), SessDur())
	putCookie(w, &cookie)
	for _, reqRole := range reqRoles {
		if reqRole == role {
			return true
		}
	}
	return false
}

/* ~ = ~ = ~ = ~ = ~ = *==[ ENCODING/DECODING ]==* ~ = ~ = ~ = ~ = ~ = */
func BaseEnc(msg string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(msg))
}

func BaseDec(data string) string {
	msg, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil {
		fmt.Printf("cook.go -> BaseDec -> base64RawURLEncoding.DecodeString() -> : %v\n", err)
		return ""
	}
	return string(msg)
}

func UrlEnc(msg string) string {
	return url.QueryEscape(msg)
}

func UrlDec(data string) string {
	msg, err := url.QueryUnescape(data)
	if err != nil {
		fmt.Printf("cook.go -> UrlDec -> url.QueryUnescape() -> : %v\n", err)
		return ""
	}
	return string(msg)
}

func Enc(data string) string {
	var eData []byte
	for _, b := range []byte(data) {
		eData = append(eData, byte(int(b)+COOKIESALT))
	}
	return BaseEnc(string(eData))
}

func Dec(eData string) string {
	eData = BaseDec(eData)
	var data []byte
	for _, b := range []byte(eData) {
		data = append(data, byte(int(b)-COOKIESALT))
	}
	return string(data)
}
