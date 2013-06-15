// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package toys

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
)

type InfoKey uint8

const (
	RequestMethod InfoKey = iota
	RequestHost
	RequestPath
	RequestQuery
	RemoteAddress
)

// Controller is the main struct of toys framework. Use it to handle the request.
type Controller struct {
	Request  *http.Request
	Response http.ResponseWriter
	inf      map[InfoKey]string
	path     string
}

// Init initial the Controller given a http.ResponseWriter and *http.Request. You must call Init
// right after create new Controller for each request.
func (c *Controller) Init(w http.ResponseWriter, r *http.Request) {
	c.Request = r
	c.Response = w
	c.inf = make(map[InfoKey]string)
	c.inf[RequestMethod] = r.Method
	c.inf[RequestHost] = r.URL.Host
	c.inf[RequestPath] = r.URL.Path
	c.inf[RequestQuery] = r.URL.RawQuery
	c.inf[RemoteAddress] = r.RemoteAddr
}

// Write writes the slice of bytes b to the web browser.
func (c *Controller) Write(b []byte) (int, error) {
	return c.Response.Write(b)
}

// Redirect send the redirect header with the url destination and the status code.
func (c *Controller) Redirect(url string, code int) {
	http.Redirect(c.Response, c.Request, url, code)
}

func (c *Controller) Info(key InfoKey) string {
	return c.inf[key]
}

// POST returns the string value for the named component of the POST or GET query. It call
// template.HTMLEscapeString for the output if filter is true.
func (c *Controller) Post(name string, filter bool) string {
	if c.Request.Method == "POST" || c.Request.Method == "PUT" {
		if filter {
			return template.HTMLEscapeString(c.Request.PostFormValue(name))
		}
		return c.Request.PostFormValue(name)
	}
	return ""
}

// Get returns the first value for the named component of the GET or HEAD query. It call
// template.HTMLEscapeString for the output if filter is true.
func (c *Controller) Get(name string, filter bool) string {
	if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
		if filter {
			return template.HTMLEscapeString(c.Request.FormValue(name))
		}
		return c.Request.FormValue(name)
	}
	return ""
}

// Cookie return the cookis value with given name. It call template.HTMLEscapeString for the output
// if filter is true.
func (c *Controller) Cookie(name string, filter bool) string {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		if filter {
			return template.HTMLEscapeString(cookie.Value)
		}
		return cookie.Value
	}
	return ""
}

// Print formats using the default formats for its operands and writes to web browser. It returns
// the number of bytes written and any write error encountered.
func (c *Controller) Print(a ...interface{}) {
	fmt.Fprint(c.Response, a...)
}

// Printf formats according to a format specifier and writes to web browser. It returns the number
// of bytes written and any write error encountered.
func (c *Controller) Printf(format string, a ...interface{}) {
	fmt.Fprintf(c.Response, format, a...)
}

// SetPath sets path of the application. For example, if you want you application handle the address
// example.com/toysapp/ then the input for SetPath should be "/toysapp/".
func (c *Controller) SetPath(p string) {
	c.path = p
}

// BasePath return the relative url base on the path seted with SetPath. BasePath will return the
// origin value if p is an absolute url.
func (c *Controller) BasePath(p string) string {
	path_url, err := url.Parse(p)
	if err != nil && path_url.IsAbs() {
		return path_url.String()
	}

	return path.Join(c.path, p)
}
