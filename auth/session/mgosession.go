// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package session

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
)

type MgoProvider struct {
	req        *http.Request
	resp       http.ResponseWriter
	cookieName string
	expiration int
	matchAddr  bool
	matchAgent bool
	collection *mgo.Collection
}

func NewMgoProvider(w http.ResponseWriter, r *http.Request,
	c *mgo.Collection) *MgoProvider {
	p := &MgoProvider{}
	p.req = r
	p.resp = w
	p.cookieName = "toysSession"
	p.expiration = 7200
	p.matchAddr = true
	p.collection = c

	return p
}

func (p *MgoProvider) SetCookieName(name string) {
	p.cookieName = name
}

func (p *MgoProvider) CookieName() string {
	return p.cookieName
}

func (p *MgoProvider) SetExpiration(exp int) {
	if exp <= 0 {
		p.expiration = 7200
		return
	}
	p.expiration = exp
}

func (p *MgoProvider) Expiration() int {
	return p.expiration
}

func (p *MgoProvider) SetMatchRemoteAddr(match bool) {
	p.matchAddr = match
}

func (p *MgoProvider) MatchRemoteAddr() bool {
	return p.matchAddr
}

func (p *MgoProvider) SetMatchUserAgent(match bool) {
	p.matchAgent = match
}

func (p *MgoProvider) MatchUserAgent() bool {
	return p.matchAgent
}

func (p *MgoProvider) load() *sessionEntry {
	entry := &sessionEntry{}
	var agent = p.req.UserAgent()
	if len(agent) > 120 {
		agent = agent[:120]
	}

	cookie, err := p.req.Cookie(p.CookieName())
	if err == nil {
		err = p.collection.FindId(cookie.Value).One(&entry)
		if err == nil {
			if (p.MatchUserAgent() && entry.UserAgent != agent) ||
				(p.MatchRemoteAddr() && entry.RemoteAddr != p.req.RemoteAddr) {
				fmt.Println("not hrer")
				p.collection.RemoveId(cookie.Value)
				return newSessionEntry(p.req.RemoteAddr, agent)
			} else {
				entry.LastActivity = time.Now()
				return entry
			}
		} else {
			fmt.Println(err.Error())
		}
	}
	entry = newSessionEntry(p.req.RemoteAddr, agent)
	http.SetCookie(p.resp, &http.Cookie{
		Name:   p.cookieName,
		Value:  entry.Id,
		MaxAge: p.Expiration(),
		Expires: entry.LastActivity.
			Add(time.Duration(p.Expiration()) * time.Second),
	})

	return entry
}

func (p *MgoProvider) setData(name string, val interface{}, flash bool) error {
	var err error
	entry := p.load()

	if flash {
		entry.FlashData[name] = val
		err = p.collection.UpdateId(entry.Id, bson.M{"$set": bson.M{
			"FlashData":    entry.FlashData,
			"LastActivity": entry.LastActivity,
		}})
	} else {
		entry.Data[name] = val
		err = p.collection.UpdateId(entry.Id, bson.M{"$set": bson.M{
			"Data":         entry.FlashData,
			"LastActivity": entry.LastActivity,
		}})
	}
	if err != nil {
		return p.collection.Insert(entry)
	}

	return nil
}

func (p *MgoProvider) Set(name string, val interface{}) error {
	return p.setData(name, val, false)
}

func (p *MgoProvider) Get(name string) interface{} {
	entry := p.load()
	fmt.Println(entry.Data[name])
	return entry.Data[name]
}

func (p *MgoProvider) GetInt(name string) int {
	n := p.Get(name)
	v, ok := n.(int)
	if !ok {
		return 0
	}
	return v
}

func (p *MgoProvider) GetBool(name string) bool {
	v, ok := p.Get(name).(bool)
	if !ok {
		return false
	}
	return v
}

func (p *MgoProvider) GetString(name string) string {
	v, ok := p.Get(name).(string)
	if !ok {
		return ""
	}
	return v
}

func (p *MgoProvider) SetFlash(name string, val interface{}) error {
	return p.setData(name, val, true)
}
