/*
* The MIT License (MIT)
*
* Copyright (c) 2017  aerth <aerth@riseup.net>
*
* Permission is hereby granted, free of charge, to any person obtaining a copy
* of this software and associated documentation files (the "Software"), to deal
* in the Software without restriction, including without limitation the rights
* to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
* copies of the Software, and to permit persons to whom the Software is
* furnished to do so, subject to the following conditions:
*
* The above copyright notice and this permission notice shall be included in all
* copies or substantial portions of the Software.
*
* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
* IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
* FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
* LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
* OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
* SOFTWARE.
 */

package tilde

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Handler is a http.Handler with log
type Handler struct {
	http.Handler
	mu       sync.Mutex
	boottime time.Time
	Log      *log.Logger
}

// NewHandler returns a http handler that serves /~tilde/ prefix
// Serves homepage, 404s, and does not serve symlinked files.
// This handler supports indexing, and will index symlink names,
// but not serve them.
func NewHandler() *Handler {
	m := new(Handler)
	os.Rename("tlde.log", "tlde.log.99")
	logfile, err := os.Create("tlde.log")
	if err != nil {
		println(err.Error())
		logfile = nil
	}

	var logger io.Writer
	if logfile == nil {
		logger = os.Stderr
	} else {
		logger = io.MultiWriter(logfile, os.Stderr)
		println("logging to:", logfile.Name())
	}
	m.boottime = time.Now()
	m.Log = log.New(logger, "[tl;de] ", log.Ltime)
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
		homehandler(w, r)
		return
	}

	if r.URL.Path == "/tlde.png" {
		png, err := Asset("tlde.png")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Write(png)
		return
	}

	if r.URL.Path == "/robots.txt" {
		robots := "User-agent: *\nDisallow: /\n"
		w.Write([]byte(robots))
		return
	}

	if r.URL.Path == "/status" {
		m.mu.Lock()
		status := time.Now().Sub(m.boottime)
		m.mu.Unlock()
		w.Write([]byte(status.String()+"\n"))
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
	if !fileisgood(abs) {
		log.Println("bad file:", abs)
		http.Error(w, "invalid file", http.StatusNotAcceptable)
		return
	} else {
	}

	// ServeFile using net/http's static file server
	m.Log.Println(r.URL.Path, "->", abs)
	http.ServeFile(w, r, abs)
	return
}

func homehandler(w http.ResponseWriter, r *http.Request) {
	_, err := os.Open(DefaultHomePageFile)
	if err == nil {
		http.ServeFile(w, r, DefaultHomePageFile)
		return
	}
	w.Write([]byte(homepage))
}
