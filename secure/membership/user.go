// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package membership

import (
	"github.com/bufio/toys/model"
	"time"
)

type User interface {
	GetId() model.Identifier
	SetId(model.Identifier) error
	GetEmail() string
	SetEmail(string)
	GetPassword() Password
	SetPassword(*Password)
	GetOldPassword() Password
	GetInfomation() Information
	SetInfomation(*Information)
	GetPrivilege() map[string]bool
	SetPrivilege(map[string]bool)
	IsApproved() bool
	Approve()
	GetConfirmCode() string
	SetConfirmCode(string)
}

type Account struct {
	Id          model.Identifier `bson:"-" datastore:"-"`
	Email       string
	OldPwd      Password
	Pwd         Password
	Info        Information
	Privilege   map[string]bool
	Approved    bool
	ConfirmCode string
}

// GetId just an virtual function, you may want to re-implement it
func (a *Account) GetId() model.Identifier {
	return a.Id
}

// SetId just an virtual function, you may want to re-implement it
func (a *Account) SetId(id model.Identifier) error {
	a.Id = id
	return nil
}

func (a *Account) GetEmail() string {
	return a.Email
}

func (a *Account) SetEmail(email string) {
	a.Email = email
}

func (a *Account) GetPassword() Password {
	return a.Pwd
}

func (a *Account) SetPassword(pwd *Password) {
	a.Pwd = *pwd
}

func (a *Account) GetOldPassword() Password {
	return a.OldPwd
}

func (a *Account) GetInfomation() Information {
	return a.Info
}

func (a *Account) SetInfomation(info *Information) {
	a.Info = *info
}

func (a *Account) GetPrivilege() map[string]bool {
	return a.Privilege
}

func (a *Account) SetPrivilege(priv map[string]bool) {
	a.Privilege = priv
}

func (a *Account) IsApproved() bool {
	return a.Approved
}

func (a *Account) Approve() {
	a.Approved = true
}

func (a *Account) GetConfirmCode() string {
	return a.ConfirmCode
}

func (a *Account) SetConfirmCode(code string) {
	a.ConfirmCode = code
}

type Password struct {
	Hashed []byte
	Salt   []byte
	InitAt time.Time
}

type Information struct {
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
