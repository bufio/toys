// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package membership

import (
	"github.com/kidstuff/toys/model"
	"time"
)

type User interface {
	GetId() model.Identifier
	SetId(model.Identifier) error
	GetEmail() string
	GetPassword() Password
	GetOldPassword() Password
	GetInfomation() UserInfo
	GetPrivilege() map[string]bool
	IsApproved() bool
	GetConfirmCodes() map[string]string
}

type Account struct {
	Id           model.Identifier `bson:"-" datastore:"-"`
	Email        string
	OldPwd       Password
	Pwd          Password
	LastActivity time.Time
	Info         UserInfo
	Privilege    map[string]bool
	Approved     bool
	ConfirmCodes map[string]string
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

func (a *Account) GetPassword() Password {
	return a.Pwd
}

func (a *Account) GetOldPassword() Password {
	return a.OldPwd
}

func (a *Account) GetInfomation() UserInfo {
	return a.Info
}

func (a *Account) GetPrivilege() map[string]bool {
	return a.Privilege
}

func (a *Account) IsApproved() bool {
	return a.Approved
}

func (a *Account) GetConfirmCodes() map[string]string {
	return a.ConfirmCodes
}

type Password struct {
	Hashed []byte
	Salt   []byte
	InitAt time.Time
}

type UserInfo struct {
	FirstName  string
	LastName   string
	MiddleName string
	NickName   string
	BirthDay   time.Time
	JoinDay    time.Time
	Address    []Address
	Phone      []string
}

type Address struct {
	Country  string
	State    string
	City     string
	District string
	Street   string
}
