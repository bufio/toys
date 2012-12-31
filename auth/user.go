// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"crypto/sha256"
	"github.com/gorilla/securecookie"
	"hash"
	"time"
)

var hashFunc func() hash.Hash = sha256.New

func SetHashFunc(f func() hash.Hash) {
	hashFunc = f
}

type Password struct {
	Hashed []byte
	Salt   []byte
	InitAt time.Time
}

type Account struct {
	Email      string
	Pwd        Password
	FirstName  string
	LastName   string
	MiddleName string
	NickName   string
	BirthDay   time.Time
	JoinDay    time.Time
	LastLogin  time.Time
	Privilege  map[string]bool
}

func NewAccount(email, password string) *Account {
	u := &Account{}

	u.privilege = make(map[string]bool)

	u.Email = email
	u.SetPassword(password)
	return u
}

func (u *Account) SetPassword(password string) {
	u.Pwd.Salt = securecookie.GenerateRandomKey(64)
	h := hashFunc()
	h.Write([]byte(password))
	h.Write(u.Pwd.Salt)
	u.Pwd.Hashed = h.Sum(nil)
	u.Pwd.InitAt = time.Now()
}
