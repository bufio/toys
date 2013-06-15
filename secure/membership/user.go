// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package membership

import (
	"github.com/bufio/toys/model"
	"time"
)

type Password struct {
	Hashed []byte
	Salt   []byte
	InitAt time.Time
}

type User struct {
	Id        model.Identifier `bson:"_id,omitempty"`
	Email     string
	OldPwd    Password
	Pwd       Password
	Info      `bson:",inline"`
	Privilege map[string]bool
	Approved  bool
	ApprCode  string
}

type Info struct {
	FirstName    string
	LastName     string
	MiddleName   string
	NickName     string
	BirthDay     time.Time
	JoinDay      time.Time
	LastActivity time.Time
	Address      []Address
	Phone        []string
}

type Address struct {
	Country  string
	State    string
	City     string
	District string
	Street   string
}
