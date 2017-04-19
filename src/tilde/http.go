package tilde

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Version string
func Version() string {
	return fmt.Sprint("tlde ", version)
}

var version = "v0.0.2"
var publichtml = "Public"

func init(){
	if os.Getenv("PUBLIC") != "" {
		publichtml = os.Getenv("PUBLIC")
	}
}

// CHMODDIR default dir permissions
var CHMODDIR = 0755 // public

// Mux is a http handler really
type Mux struct {
	Log *log.Logger
}

// Handler returns a http handler that serves /~tilde/
func Handler() *Mux {
	m := new(Mux)
	os.MkdirAll("logs", os.FileMode(CHMODDIR))
	logfile, err := ioutil.TempFile("logs", "tlde")
	var mw io.Writer
	if err != nil {
		logfile = os.Stderr
		mw = os.Stderr
	} else {
		mw = io.MultiWriter(logfile, os.Stderr)
	}
	println("logging to", logfile.Name())
	m.Log = log.New(mw, "[tl;de] ", log.Lshortfile)
	return m
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "tlde/"+strings.TrimPrefix(version, "v"))
	m.Log.Println(r.Method, r.URL.Path)
	if r.URL.Path == "/" {
		w.Write([]byte(homepage))
		return
	}
	if !strings.HasPrefix(r.URL.Path, "/~") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
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

	// public folder
	dir := fmt.Sprintf("/home/%s/%s", u, publichtml)

	// let net/http FileServer handle the rest
	handler := http.StripPrefix("/~"+u+"/", http.FileServer(http.Dir(dir)))
	handler.ServeHTTP(w, r)
}

var homepage = ``+
`<!DOCTYPE html>
<html lang="en">
  <head>
    <title>&lrm;</title>
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