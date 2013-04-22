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

type Controller struct {
	req   *http.Request
	respw http.ResponseWriter
	inf   map[string]string
	path  string
}

func (c *Controller) Init(w http.ResponseWriter, r *http.Request) {
	c.req = r
	c.respw = w
	c.inf = make(map[string]string)
	c.inf["req_method"] = r.Method
	c.inf["req_host"] = r.URL.Host
	c.inf["req_path"] = r.URL.Path
	c.inf["req_query"] = r.URL.RawQuery
	c.inf["remote_addr"] = r.RemoteAddr
}

func (c *Controller) Write(b []byte) (int, error) {
	return c.respw.Write(b)
}

func (c *Controller) Request() *http.Request {
	return c.req
}

func (c *Controller) Redirect(url string, code int) {
	http.Redirect(c.respw, c.req, url, code)
}

func (c *Controller) Info(key string) string {
	return c.inf[key]
}

func (c *Controller) Post(name string, filter bool) string {
	if c.req.Method == "POST" {
		if filter {
			return template.HTMLEscapeString(c.req.FormValue(name))
		}
		return c.req.FormValue(name)
	}
	return ""
}

func (c *Controller) Get(name string, filter bool) string {
	if c.req.Method == "GET" || c.req.Method == "HEAD" {
		if filter {
			return template.HTMLEscapeString(c.req.FormValue(name))
		}
		return c.req.FormValue(name)
	}
	return ""
}

func (c *Controller) Cookie(name string, filter bool) string {
	cookie, err := c.req.Cookie(name)
	if err != nil {
		if filter {
			return template.HTMLEscapeString(cookie.Value)
		}
		return cookie.Value
	}
	return ""
}

func (c *Controller) Print(a ...interface{}) {
	fmt.Fprint(c.respw, a...)
}

func (c *Controller) Printf(format string, a ...interface{}) {
	fmt.Fprintf(c.respw, format, a...)
}

func (c *Controller) SetPath(p string) {
	c.path = p
}

func (c *Controller) BasePath(p string) string {
	path_url, err := url.Parse(p)
	if err != nil && path_url.IsAbs() {
		return path_url.String()
	}

	return path.Join(c.path, p)
}
