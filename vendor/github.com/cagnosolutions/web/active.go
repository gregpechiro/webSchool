package web

import (
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ActivityManager = NewActivityMan()

var AMANAGER = false

type Activity struct {
	Time      int64
	Ip        string
	Actions   Queue
	UserAgent string
	Session   session
}

func NewActivity() *Activity {
	return &Activity{
		Actions: NewSQueue(10),
	}
}

func (a *Activity) Update(r *http.Request) {
	t := time.Now()
	a.Time = t.UnixNano()
	a.Ip = r.Header.Get("X-Real-IP")
	a.Session = GetAllSess(r)
	if a.Ip == "" {
		a.Ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	a.Actions.Put(NewAction(r, t))
}

type ActivityOrder []*Activity

func (s ActivityOrder) Len() int {
	return len(s)
}

func (s ActivityOrder) Less(i, j int) bool {
	return s[i].Time < s[j].Time
}

func (s ActivityOrder) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ActivityMan struct {
	Activities map[string]*Activity
	sync.RWMutex
}

func NewActivityMan() *ActivityMan {
	return &ActivityMan{
		Activities: make(map[string]*Activity),
	}
}

func (s *ActivityMan) New(r *http.Request) string {
	t := time.Now()
	rnd := rand.New(rand.NewSource(t.UnixNano()))
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	actId := ip + "::" + strconv.Itoa(rnd.Intn(1000))
	act := NewActivity()
	act.Time = t.UnixNano()
	act.Ip = ip
	act.Actions.Put(NewAction(r, t))
	act.UserAgent = r.UserAgent()
	act.Session = GetAllSess(r)
	ActivityManager.Put(actId, act)
	return actId
}

func (s *ActivityMan) Put(key string, act *Activity) {
	s.Lock()
	s.Activities[key] = act
	s.Unlock()
}

func (s *ActivityMan) Get(key string) (*Activity, bool) {
	s.RLock()
	act, ok := s.Activities[key]
	s.RUnlock()
	return act, ok
}

func (s *ActivityMan) All() []*Activity {
	var acts ActivityOrder
	s.RLock()
	for _, ss := range s.Activities {
		acts = append(acts, ss)
	}
	s.RUnlock()
	sort.Stable(sort.Reverse(acts))
	return acts
}

type Action struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Time   string `json:"time"`
}

func NewAction(r *http.Request, t time.Time) Action {
	return Action{
		Method: r.Method,
		Path:   r.URL.Path,
		Time:   t.Format("1/2/2006 3:04:05 PM"),
	}
}

type Queue struct {
	Q   []interface{}
	Max int
}

func NewSQueue(max int) Queue {
	return Queue{
		Q:   make([]interface{}, 0, 0),
		Max: max,
	}
}

func (q *Queue) Put(v interface{}) {
	if len(q.Q) >= q.Max {
		q.Q = append(q.Q[1:], v)
		return
	}
	q.Q = append(q.Q, v)
}

func (q *Queue) All() []interface{} {
	var nq []interface{}
	for i := len(q.Q) - 1; i > -1; i-- {
		nq = append(nq, q.Q[i])
	}
	return nq
}

func (q *Queue) ToJson() string {
	b, err := json.Marshal(q.Q)
	if err != nil {
		return ""
	}
	return string(b)
}

func (q *Queue) Recent() interface{} {
	if len(q.Q) > 0 {
		return q.Q[len(q.Q)-1]
	}
	return nil
}

func handleActivity(w http.ResponseWriter, r *http.Request) {
	//t, err := template.ParseFiles("activity.tmpl")

	t, err := template.New("activityPage").Parse(activityPage)
	if err != nil {
		log.Printf("package web >> active.go >> handleActivity() >> template.ParseFiles() >> %v\n\n", err)
		return
	}
	t.Execute(w, Model{
		"uaJs":       UAJS,
		"activities": ActivityManager.All(),
	})
}
