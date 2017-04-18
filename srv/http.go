package srv

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var publichtml = "Public" // consider public_html instead of Public
// CHMODDIR default dir permissions
var CHMODDIR = 0755       // public
// Mux is a router
type Mux struct {
	Log *log.Logger
}

// NewMux ..
func NewMux() *Mux {
	m := new(Mux)
	os.MkdirAll("logs", os.FileMode(CHMODDIR))
	logfile, err := ioutil.TempFile("logs", "tlde")
	if err != nil {
		logfile = os.Stderr
	}
	m.Log = log.New(logfile, "[tl;de] ", log.Lshortfile)
	return m
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Log.Println(r.URL.Path)
	if !strings.HasPrefix(r.URL.Path, "/~") {
		buf := new(bytes.Buffer)
		buf.WriteString("#!/bin/sh\n")
		buf.WriteString("echo wtf r u doing sheller\n")
		buf.WriteString("echo danger zone\n")
		w.Write(buf.Bytes())
		return
	}
	var u string
	u = strings.Split(strings.TrimPrefix(r.URL.Path, "/~"), "/")[0]
	m.Log.Println(u)
	m.Log.Println(r.URL.Path)

	if u != "" {
	dir := fmt.Sprintf("/home/%s/%s", u, publichtml)
	m.Log.Println("Serving Dir:", dir)
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/~"+u)
	m.Log.Println(r.URL.Path)
	ha := http.FileServer(http.Dir(dir))
	
	ha.ServeHTTP(w,r)
	}
}
