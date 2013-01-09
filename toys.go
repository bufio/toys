// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package toys

import (
	"html/template"
	"net/http"
)

type Controller struct {
	req   *http.Request
	respw http.ResponseWriter
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
		return cookie.Value
	}
	return ""
}
