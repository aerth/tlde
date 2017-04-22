package tilde

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Version string
func Version() string {
	return fmt.Sprint("tlde ", version)
}

var version = "v0.0.3"

// PublicDir such as Public or public_html
var PublicDir = "Public"

// FormatPath is how we find user home dir
// first var is username, second var is publichtml
// for /home/user/[publichtml]
// consider /path/to/home/%s/%s
var FormatPath = "/home/%s/%s"

//var formatpath = "/usr/home/%s/%s"

func init() {
	if os.Getenv("PUBLIC") != "" {
		PublicDir = os.Getenv("PUBLIC")
	}
	if os.Getenv("FORMATPATH") != "" {
		FormatPath = os.Getenv("FORMATPATH")
	}
}

// CHMODDIR default dir permissions
var CHMODDIR = 0755 // public

// Handler handles http requests
type Handler struct {
	Log *log.Logger
}

// NewHandler returns a http handler that serves /~tilde/ prefix
// Serves homepage, 404s, and does not serve symlinked files.
// This handler supports indexing, and will index symlink names,
// but not serve them.
func NewHandler() *Handler {
	m := new(Handler)
	logfile, err := os.Create("tlde.log")
	if err != nil {
		println(err.Error())
		logfile = nil
	}

	var mw io.Writer
	if logfile == nil {
		mw = os.Stderr
	} else {
		mw = io.MultiWriter(logfile, os.Stderr)
		println("logging to:", logfile.Name())
	}

	m.Log = log.New(mw, "[tl;de] ", log.Ltime)
	return m
}

func (m *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "tlde/"+strings.TrimPrefix(version, "v"))
	m.Log.Println(r.Method, r.URL.Path)

	// GET only
	if r.Method != "GET" {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Handle "/" (homepage)
	if r.URL.Path == "/" {
		w.Write([]byte(homepage))
		return
	}

	// Redirect non-tilde URLs back to homepage
	if !strings.HasPrefix(r.URL.Path, "/~") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// extract user name from URL
	var u string
	u = strings.Split(strings.TrimPrefix(r.URL.Path, "/~"), "/")[0]
	if u == "" {
		http.NotFound(w, r)
		return
	}

	// if no slash+tilde+user+slash, redirect to slash+tilde+user+slash
	if !strings.HasPrefix(r.URL.Path, "/~"+u+"/") {
		new := strings.Replace(r.URL.Path, "/~"+u, "/~"+u+"/", 1)
		http.Redirect(w, r, new, http.StatusFound)
		return
	}

	// requested file path (ex: /doc/readme.md from /~user/doc/readme.md)
	reqfile := strings.TrimPrefix(r.URL.Path, "/~"+u+"/")

	// prevent rare case of PublicDir being empty string, thereby serving '/home/user//'
	if PublicDir == "" {
		http.Error(w, "system failure", http.StatusInternalServerError)
		return
	}

	// format requested file into filesystem path
	// (ex: /path/to/file.txt becomes /home/user/Public/path/to/file.txt)
	reqfile = filepath.Join(fmt.Sprintf(FormatPath, u, PublicDir), reqfile)

	// absolute filesystem path of requested file
	abs, err := filepath.Abs(reqfile)
	if err != nil {
		http.Error(w, "bad path", http.StatusUnauthorized)
		return
	}

	// 404
	if _, err := os.Stat(abs); err != nil {
		m.Log.Println(err)
		http.NotFound(w, r)
		return
	}

	// is symlink
	if !m.isgood(abs) {
		http.Error(w, "invalid file", http.StatusForbidden)
		return
	}

	// ServeFile using net/http's static file server
	m.Log.Println(r.URL.Path, "->", abs)
	http.ServeFile(w, r, abs)
	return
}

// homepage html
var homepage = `` +
	`<!DOCTYPE html>
<html lang="en">
  <head>
    <title>#!</title>
    <meta name="viewport" content="initial-scale = 1, maximum-scale=1, user-scalable = 0"/>
    <meta name="apple-mobile-web-app-capable" content="yes"/>
    <meta name="apple-mobile-web-app-status-bar-style" content="black"/>
    <meta name="HandheldFriendly" content="true"/>
    <meta name="MobileOptimized" content="320"/>
    <link href="https://fonts.googleapis.com/css?family=Montserrat" rel="stylesheet" type="text/css">
    <link href="http://hashbang.sh/assets/local.css" rel="stylesheet" type="text/css">
  </head>
  <body>
    <script src="http://hashbang.sh/assets/icon.js"></script>
    <h1>#!</h1>
    <a href="view-source:https://hashbang.sh">
      <code>sh <(curl hashbang.sh | gpg)</code>
    </a>
  </body>
</html>
`

// returns false if symlink
// comparing absolute vs resolved path is quick and effective
func (h *Handler) isgood(abs string) bool {
	realpath, err := filepath.EvalSymlinks(abs)
	if err != nil {
		h.Log.Println(err)
		return false
	}
	if realpath != abs {
		h.Log.Println("not serving symlink to", realpath, "from", abs)
		return false
	}
	
	return true
}
