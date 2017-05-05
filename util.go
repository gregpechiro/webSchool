package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func encodeFile(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func setFileType(name string) string {
	switch filepath.Ext(name) {
	case ".html":
		return "html"
	case ".js":
		return "javascript"
	case ".css":
		return "css"
	case ".json":
		return "json"
	case ".jpg":
		return "image"
	case ".png":
		return "image"
	}
	return "default"
}

func defaultUsers() {
	admin := User{
		Id:        "0",
		Role:      "ADMIN",
		FirstName: "Admin",
		LastName:  "Temporary",
		Email:     "admin@temp.com",
		Password:  "admin",
		Active:    true,
	}

	db.Set("user", "0", admin)

	fmt.Printf("\nTemporary admin credentials:\n\n\tEmail:\t\t%s\n\tPassword:\t%s\n\n", admin.Email, admin.Password)
}

func ajaxResponse(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, msg)
}

func genId() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
}

func DirStats(path string) (int64, int, error) {
	var size int64
	var files int
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files++
			size += info.Size()
		}
		return err
	})
	return size, files, err
}

func IsEmptyDir(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return true
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err != nil {
		return true
	}
	return false
}

func PrettySize(size int64) string {
	c := 0
	var sizef float64 = float64(size)
	for sizef > 1024 {
		sizef = sizef / 1024
		c++
	}
	ind := ""
	switch c {
	case 0:
		ind = "B"
	case 1:
		ind = "KB"
	case 2:
		ind = "MB"
	case 3:
		ind = "GB"
	}
	return fmt.Sprintf("%.1f %s", sizef, ind)
}

func pretty(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.ToUpper(string(s[0])) + s[1:]
	for i := 0; i < len(s); i++ {
		if s[i] == byte(' ') {
			s = s[:i+1] + strings.ToUpper(string(s[i+1])) + s[i+2:]
		}
	}
	return s
}

func PrettyDateTime(ts int64) string {
	if ts == 0 {
		return ""
	}
	t := time.Unix(ts, 0)
	return t.Format("1/2/2006 3:04 PM")
}

var themes = []string{
	"ambiance",
	"chaos",
	"chrome",
	"clouds",
	"clouds_midnight",
	"cobalt",
	"crimson_editor",
	"dawn",
	"dreamweaver",
	"eclipse",
	"github",
	//"gruvbox",
	"idle_fingers",
	"iplastic",
	"katzenmilch",
	"kr_theme",
	"kuroir",
	"merbivore",
	"merbivore_soft",
	"mono_industrial",
	"monokai",
	"pastel_on_dark",
	"solarized_dark",
	"solarized_light",
	"sqlserver",
	"terminal",
	"textmate",
	"tomorrow",
	"tomorrow_night",
	"tomorrow_night_blue",
	"tomorrow_night_bright",
	"tomorrow_night_eighties",
	"twilight",
	"vibrant_ink",
	"xcode",
}
