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

// DefaultHomePageFile is where tlde looks for the home page.
// If it doesn't exist, homepage (defined below) will be served in its place.
// This only gets seen for 'GET /' requests
var DefaultHomePageFile = "index.html"

// homepage html
var homepage = `` +
	`<!DOCTYPE html>
<html lang="en">
  <head>
    <title>tilde server</title>
    <meta name="viewport" content="initial-scale = 1, maximum-scale=1, user-scalable = 0"/>
    <meta name="apple-mobile-web-app-capable" content="yes"/>
    <meta name="apple-mobile-web-app-status-bar-style" content="black"/>
    <meta name="HandheldFriendly" content="true"/>
    <meta name="MobileOptimized" content="320"/>
  </head>
  <body style="background-color: black; color: green;" text-align: center;">
<style>
html {
  background:black;
  color:black;
}
body {
  display:block;
  position:fixed;
  text-align: center;
  width:320px;
  height:320px;
  overflow:hidden;
}
img {
	max-width: 50vw;
	height: auto;
}
p {
	font-size: larger;
}
</style>
	<div style="width: 100vw; text-align: center; margin-top: 20vh;">
	<a href="https://github.com/aerth/tlde">
	<img src="/tlde.png" alt="tlde server https://github.com/aerth/tlde">
	</a>

	<p>invalid path: try <code>/~username/</code> </p>
  </body>
</html>
`
