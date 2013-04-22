// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package sessions

import (
	"encoding/base64"
	"github.com/openvn/toys/secure"
	"time"
)

type sessionEntry struct {
	Id           string `bson:"_id"`
	RemoteAddr   string
	UserAgent    string
	LastActivity time.Time
	Data         map[string]interface{}
	FlashData    map[string]interface{}
}

func newSessionEntry(addr, agent string) *sessionEntry {
	s := &sessionEntry{}
	s.Id = base64.URLEncoding.EncodeToString(secure.RandomToken(32))
	s.RemoteAddr = addr
	s.UserAgent = agent
	s.LastActivity = time.Now()
	s.Data = make(map[string]interface{})
	s.FlashData = make(map[string]interface{})
	return s
}

type Provider interface {
	SetCookieName(name string)
	CookieName() string
	SetExpiration(exp int)
	Expiration() int
	SetMatchRemoteAddr(match bool)
	MatchRemoteAddr() bool
	SetMatchUserAgent(match bool)
	MatchUserAgent() bool

	Set(name string, val interface{}) error
	Get(name string) interface{}
	GetInt(name string) int
	GetBool(name string) bool
	GetString(name string) string
	Delete(name ...string) error
	DeleteAll(flash bool) error

	SetFlash(name string, val interface{}) error
	GetFlash(name string) interface{}
	GetFlashInt(name string) int
	GetFlashBool(name string) bool
	GetFlashString(name string) string

	Destroy() error
}
