// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"encoding/base64"
	"errors"
	"github.com/openvn/toys/secure"
	"hash"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strings"
	"time"
)

const (
	CookieName string = "toysAuth"
)

var (
	ErrInvalidId       error = errors.New("auth: invalid id, the id must be in bson.ObjectId format")
	ErrInvalidEmail    error = errors.New("auth: invalid email address")
	ErrDuplicateEmail  error = errors.New("auth: duplicate email address")
	ErrInvalidPassword error = errors.New("auth: invalid password")
)

type Config struct {
	Key   string `bson:"_id"`
	Value interface{}
}

type rememberInfo struct {
	Id    bson.ObjectId `bson:"_id"`
	Token string
	Exp   time.Time
}

type Authenticater interface {
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
	AddUserInfo(email, pwd string, info Info, pri map[string]bool, notif, app bool) error
	// DeleteUserByEmail deletes an user from database base on the given id;
	// It returns an error describes the first issue encountered, if any.
	DeleteUser(id string) error
	// GetUser gets the infomations and update the LastActivity of the current
	// logged user;
	// It returns an error describes the first issue encountered, if any.
	GetUser() (*User, error)
	// FindUser finds the user with the given id;
	// Its returns an ErrNotFound if the user's id was not found.
	FindUser(id string) (*User, error)
	// FindUserByEmail like FindUser but receive an email
	FindUserByEmail(email string) (*User, error)
	// FindAllUser finds and return a slice of user.
	// offsetId, limit define which sub-sequence of matching users to return.
	// Limit take an number of user per page; offsetId take the Id of the last
	// user of the previous page.
	FindAllUser(offsetId string, limit int) ([]*User, error)
	// FindAllUserOline finds and return a slice of current logged user.
	// See FindAllUser for the usage.
	FindUserOnline(offsetId string, limit int) ([]*User, error)
	// CountUserOnline counts the number of user current logged.
	// It counts the user that LastActivity+OnlineThreshold<Now.
	CountUserOnline() int
	// ValidateUser validate user email and password.
	// It returns the user infomations if the email and password is correct.
	ValidateUser(email string, password string) (*User, bool)
	// LogginUser logs user in and set the session "user_email" with value is
	// the user email string. Remember take a number of second to keep the user
	// loggin state.
	// Developer must call LogginUser before send any output to browser.
	LogginUser(email string, remember int)
}

type AuthDBCtx struct {
	req          *http.Request
	respw        http.ResponseWriter
	notifer      Notificater
	fmtChecker   FormatChecker
	pwdHash      hash.Hash
	userColl     *mgo.Collection
	confgColl    *mgo.Collection
	rememberColl *mgo.Collection
}

func NewAuthDBCtx() Authenticater {
	a := &AuthDBCtx{}
	return a
}

func (a *AuthDBCtx) SetNotificater(n Notificater) {
	a.notifer = n
}

func (a *AuthDBCtx) SetHashFunc(h hash.Hash) {
	a.pwdHash = h
}

func (a *AuthDBCtx) SetFormatChecker(c FormatChecker) {
	a.fmtChecker = c
}

func (a *AuthDBCtx) createUser(email, password string, app bool) (*User, error) {
	if !a.fmtChecker.EmailValidate(email) {
		return nil, ErrInvalidEmail
	}
	if !a.fmtChecker.PasswordValidate(password) {
		return nil, ErrInvalidPassword
	}

	u := &User{}
	u.Email = email
	u.Pwd.InitAt = time.Now()
	u.Pwd.Salt = secure.RandomToken(32)
	a.pwdHash.Write([]byte(password))
	a.pwdHash.Write(u.Pwd.Salt)
	u.Pwd.Hashed = a.pwdHash.Sum(nil)
	a.pwdHash.Reset()

	u.Approved = app
	return u, nil
}

func (a *AuthDBCtx) insertUser(u *User, notif, app bool) error {
	err := a.userColl.Insert(u)
	if err != nil {
		if mgo.IsDup(err) {
			return ErrDuplicateEmail
		}
		return err
	}

	if notif {
		return a.notifer.AccountAdded(u.Email, app)
	}
	return nil
}

func (a *AuthDBCtx) AddUser(email, password string, notif, app bool) error {
	u, err := a.createUser(email, password, app)
	if err != nil {
		return err
	}

	return a.insertUser(u, notif, app)
}

func (a *AuthDBCtx) AddUserInfo(email, password string, info Info,
	pri map[string]bool, notif, app bool) error {
	u, err := a.createUser(email, password, app)
	if err != nil {
		return err
	}

	u.Info = info
	u.Privilege = pri

	return a.insertUser(u, notif, app)
}

func validateObjectHex(id string) bool {
	if len(id) == 12 {
		return true
	}
	return false
}

func (a *AuthDBCtx) DeleteUser(id string) error {
	if validateObjectHex(id) {
		return a.userColl.RemoveId(bson.ObjectIdHex(id))
	}
	return ErrInvalidId
}

func (a *AuthDBCtx) GetUser() (*User, error) {
	//check for remember cookie
	cookie, err := a.req.Cookie(CookieName)
	if err == nil {
		pos := strings.Index(cookie.Value, "|")
		id := cookie.Value[:pos]
		token := cookie.Value[pos+1:]
		if validateObjectHex(id) {
			r := rememberInfo{}
			oid := bson.ObjectIdHex(id)
			err = a.rememberColl.FindId(oid).One(&r)
			if err == nil {
				if token == r.Token {
					user := User{}
					err = a.userColl.FindId(oid).One(&user)
					if err == nil {
						token = base64.URLEncoding.EncodeToString(secure.RandomToken(128))
						http.SetCookie(a.respw, &http.Cookie{
							Name:  CookieName,
							Value: id + "|" + token,
						})
						err = a.rememberColl.UpdateId(oid, bson.M{
							"$set": bson.M{"token": token},
						})
						if err == nil {
							return &user, nil
						} else {
							http.SetCookie(a.respw, &http.Cookie{
								Name:   CookieName,
								MaxAge: -1,
							})
							a.rememberColl.RemoveId(oid)
						}
					}
				}
			}
		}
	}
	//check for session
	//not logged-in
	return nil, errors.New("auth: not logged-in")
}

func (a *AuthDBCtx) FindUser(id string) (*User, error) {
	u := &User{}
	return u, nil
}

func (a *AuthDBCtx) FindUserByEmail(email string) (*User, error) {
	u := &User{}
	return u, nil
}

func (a *AuthDBCtx) FindAllUser(offsetKey string, limit int) ([]*User, error) {
	u := []*User{}
	return u, nil
}

func (a *AuthDBCtx) FindUserOnline(offsetKey string, limit int) ([]*User, error) {
	u := []*User{}
	return u, nil
}

func (a *AuthDBCtx) CountUserOnline() int {
	return 0
}

func (a *AuthDBCtx) ValidateUser(email string, password string) (*User, bool) {
	u := &User{}
	return u, false
}

func (a *AuthDBCtx) LogginUser(email string, remember int) {

}
