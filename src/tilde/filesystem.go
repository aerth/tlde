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
	"os"
	"path/filepath"
)

// PublicDir such as Public or public_html
var PublicDir = "Public"

// FormatPath is how we find user home dir
// first var is username, second var is publichtml
// for /home/user/[publichtml]
// consider /path/to/home/%s/%s
var FormatPath = "/home/%s/%s"

//var FormatPath = "/usr/home/%s/%s"

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

// returns false if symlink
// comparing absolute vs resolved path is quick and effective
func fileisgood(abs string) bool {
	realpath, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return false
	}
	return realpath == abs
}
