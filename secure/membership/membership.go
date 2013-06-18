// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package membership provide a interface to easy authorization for your web application.
*/
package membership

import (
	"errors"
	"github.com/bufio/toys/model"
	"hash"
)

var (
	ErrInvalidId       error = errors.New("auth: invalid id")
	ErrInvalidEmail    error = errors.New("auth: invalid email address")
	ErrDuplicateEmail  error = errors.New("auth: duplicate email address")
	ErrInvalidPassword error = errors.New("auth: invalid password")
)

type Authenticater interface {
	SetPath(p string)
	SetDomain(d string)
	// SetOnlineThreshold sets the online threshold time, if t <= 0, the loggin
	// state will last until the session expired.
	SetOnlineThreshold(t int)
	// SetHashFunc sets the hash.Hash which will be use for password hasing
	SetHashFunc(h hash.Hash)
	// SetNotificater sets the Notificater which will be use for sending
	// notification to user when account added or password changed.
	SetNotificater(n Notificater)
	// SetFormatChecker sets a FormatChecker for validate email/password
	SetFormatChecker(c FormatChecker)
	// AddUser adds an user to database with email and password;
	// if notif is true, a NewAccount notification will be send to user by the
	// Notificater. If app is false, the user is waiting to be approved.
	// It returns an error describes the first issue encountered, if any.
	AddUser(email, pwd string, notif, app bool) error
	// AddUserInfo adds an user to database;
	// if notif is true, a NewAccount notification will be send to user by the
	// Notificater. If app is false, the user is waiting to be approved.
	// It returns an error describes the first issue encountered, if any.
	AddUserInfo(email, pwd string, info *Information, pri map[string]bool, notif, app bool) error
	// DeleteUserByEmail deletes an user from database base on the given id;
	// It returns an error describes the first issue encountered, if any.
	DeleteUser(id model.Identifier) error
	// GetUser gets the infomations and update the LastActivity of the current
	// logged user;
	// It returns an error describes the first issue encountered, if any.
	GetUser() (User, error)
	// FindUser finds the user with the given id;
	// Its returns an ErrNotFound if the user's id was not found.
	FindUser(id model.Identifier) (User, error)
	// FindUserByEmail like FindUser but receive an email
	FindUserByEmail(email string) (User, error)
	// FindAllUser finds and return a slice of user.
	// offsetId, limit define which sub-sequence of matching users to return.
	// Limit take an number of user per page; offsetId take the Id of the last
	// user of the previous page.
	FindAllUser(offsetId model.Identifier, limit int) ([]User, error)
	// FindAllUserOline finds and return a slice of current logged user.
	// See FindAllUser for the usage.
	FindUserOnline(offsetId model.Identifier, limit int) ([]User, error)
	// CountUserOnline counts the number of user current logged.
	// It counts the user that LastActivity+OnlineThreshold<Now.
	CountUserOnline() int
	// ValidateUser validate user email and password.
	// It returns the user infomations if the email and password is correct.
	ValidateUser(email string, password string) (User, error)
	// LogginUser logs user in by using a session that store user id.
	// Remember take a number of second to keep the user loggin state.
	// Developer must call LogginUser before send any output to browser.
	LogginUser(id model.Identifier, remember int) error
}

// type Config struct {
// 	Key   string `bson:"_id"`
// 	Value interface{}
// }
