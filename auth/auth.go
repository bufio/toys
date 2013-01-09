// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"labix.org/v2/mgo"
	"time"
)

type Authenticater interface {
	SetOnlineTime(d time.Duration)
	GetOnlineTime() time.Duration
	AddUser(u User) error
	DeleteUser(email string) error
	GetUser() (User, error)
	FindUser(id string) (User, error)
	FindUserByEmail(email string) (User, error)
	FindAllUser(offsetKey string, limit int) ([]User, error)
	FindUserOnline(offsetKey string, limit int) ([]User, error)
	CountUserOnline() int
	ValidateUser(email string, password string) (User, bool)
	LogginUser(u User, remember int)
}

type AuthDBCtx struct {
	sess        *mgo.Session
	confgCol    *mgo.Collection
	rememberCol *mgo.Collection
}

func (a *AuthDBCtx) Close() {
	a.sess.Close()
}

func (a *AuthDBCtx) SetOnlineTime(d time.Duration) {
}
