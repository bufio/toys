// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package toys

import (
	"fmt"
	"html/template"
	"net/http"
)

type Controller struct {
	req   *http.Request
	respw http.ResponseWriter
}

func (c *Controller) Init(w http.ResponseWriter, r *http.Request) {
	c.req = r
	c.respw = w
}

func (c *Controller) Write(b []byte) (int, error) {
	return c.respw.Write(b)
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
