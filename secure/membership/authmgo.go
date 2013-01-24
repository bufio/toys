// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package membership

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/openvn/toys/secure"
	"github.com/openvn/toys/secure/membership/session"
	"hash"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strings"
	"time"
)

type AuthMongoDBCtx struct {
	threshold    time.Duration
	sess         session.Provider
	req          *http.Request
	respw        http.ResponseWriter
	notifer      Notificater
	fmtChecker   FormatChecker
	pwdHash      hash.Hash
	userColl     *mgo.Collection
	rememberColl *mgo.Collection
	cookieName   string
	sessionName  string
}

func NewAuthDBCtx(w http.ResponseWriter, r *http.Request, sess session.Provider,
	userColl, rememberColl *mgo.Collection) Authenticater {
	a := &AuthMongoDBCtx{}
	a.respw = w
	a.req = r
	a.sess = sess
	a.userColl = userColl
	a.rememberColl = rememberColl
	a.cookieName = "toysAuthCookie"
	a.sessionName = "toysAuthSession"
	a.fmtChecker, _ = NewSimpleChecker(8)
	a.notifer = NewSimpleNotificater()
	a.pwdHash = sha256.New()
	a.threshold = 900 * time.Second
	return a
}

func (a *AuthMongoDBCtx) SetOnlineThreshold(t int) {
	if t > 0 {
		a.threshold = time.Duration(t) * time.Second
	}
}

func (a *AuthMongoDBCtx) SetNotificater(n Notificater) {
	a.notifer = n
}

func (a *AuthMongoDBCtx) SetHashFunc(h hash.Hash) {
	a.pwdHash = h
}

func (a *AuthMongoDBCtx) SetFormatChecker(c FormatChecker) {
	a.fmtChecker = c
}

func (a *AuthMongoDBCtx) createUser(email, password string, app bool) (*User, error) {
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

func (a *AuthMongoDBCtx) insertUser(u *User, notif, app bool) error {
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

func (a *AuthMongoDBCtx) AddUser(email, password string, notif, app bool) error {
	u, err := a.createUser(email, password, app)
	if err != nil {
		return err
	}

	return a.insertUser(u, notif, app)
}

func (a *AuthMongoDBCtx) AddUserInfo(email, password string, info Info,
	pri map[string]bool, notif, app bool) error {
	u, err := a.createUser(email, password, app)
	if err != nil {
		return err
	}

	u.Info = info
	u.Privilege = pri

	return a.insertUser(u, notif, app)
}

func (a *AuthMongoDBCtx) DeleteUser(id string) error {
	if bson.IsObjectIdHex(id) {
		return a.userColl.RemoveId(bson.ObjectIdHex(id))
	}
	return ErrInvalidId
}

func (a *AuthMongoDBCtx) GetUser() (*User, error) {
	//check for remember cookie
	cookie, err := a.req.Cookie(a.cookieName)
	if err == nil {
		//read and parse cookie
		pos := strings.Index(cookie.Value, "|")
		id := cookie.Value[:pos]
		token := cookie.Value[pos+1:]
		if bson.IsObjectIdHex(id) {
			r := rememberInfo{}
			oid := bson.ObjectIdHex(id)
			//validate
			err = a.rememberColl.FindId(oid).One(&r)
			if err == nil {
				if token == r.Token {
					if r.Exp.After(time.Now()) {
						//delete expried auth
						goto DelCookie
					}
					user := User{}
					err = a.userColl.FindId(oid).One(&user)
					if err == nil {
						//re-generate token
						token = base64.URLEncoding.EncodeToString(secure.RandomToken(128))
						http.SetCookie(a.respw, &http.Cookie{
							Name:    a.cookieName,
							Value:   id + "|" + token,
							Expires: r.Exp,
						})
						err = a.rememberColl.UpdateId(oid, bson.M{
							"$set": bson.M{"token": token},
						})
						if err == nil {
							return &user, nil
						}
					}
				}
			}
			a.rememberColl.RemoveId(oid)
		}
	DelCookie:
		http.SetCookie(a.respw, &http.Cookie{
			Name:   a.cookieName,
			MaxAge: -1,
		})
	}
	//check for session
	inf, ok := a.sess.Get(a.sessionName).(sessionInfo)
	println("get session")
	if ok {
		println("ok")
		if inf.At.Add(a.threshold).Before(time.Now()) {
			println("valid")
			user := User{}
			err = a.userColl.FindId(inf.Id).One(&user)
			if err == nil {
				return &user, nil
			}
		} else {
			println("no valid")
			a.sess.Delete(a.sessionName)
		}
	}
	//not logged-in
	return nil, errors.New("auth: not logged-in")
}

func (a *AuthMongoDBCtx) FindUser(id string) (*User, error) {
	u := &User{}
	return u, nil
}

func (a *AuthMongoDBCtx) FindUserByEmail(email string) (*User, error) {
	u := &User{}
	return u, nil
}

func (a *AuthMongoDBCtx) FindAllUser(offsetKey string, limit int) ([]*User, error) {
	u := []*User{}
	return u, nil
}

func (a *AuthMongoDBCtx) FindUserOnline(offsetKey string, limit int) ([]*User, error) {
	u := []*User{}
	return u, nil
}

func (a *AuthMongoDBCtx) CountUserOnline() int {
	return 0
}

func (a *AuthMongoDBCtx) ValidateUser(email string, password string) (*User, error) {
	u := &User{}
	err := a.userColl.Find(bson.M{"email": email}).One(&u)
	if err != nil {
		return nil, err
	}
	a.pwdHash.Write([]byte(password))
	a.pwdHash.Write(u.Pwd.Salt)
	hashed := a.pwdHash.Sum(nil)
	a.pwdHash.Reset()
	if bytes.Compare(u.Pwd.Hashed, hashed) != 0 {
		return nil, err
	}
	return u, nil
}

func (a *AuthMongoDBCtx) LogginUser(id string, remember int) error {
	if !bson.IsObjectIdHex(id) {
		return ErrInvalidId
	}
	oid := bson.ObjectIdHex(id)
	if remember > 0 {
		//use cookie a rememberColl
		r := rememberInfo{}
		r.Id = oid
		r.Exp = time.Now().Add(time.Duration(remember) * time.Second)
		r.Token = base64.URLEncoding.EncodeToString(secure.RandomToken(128))
		http.SetCookie(a.respw, &http.Cookie{
			Name:    a.cookieName,
			Value:   id + "|" + r.Token,
			Expires: r.Exp,
		})
		return a.rememberColl.Insert(&r)
	} else {
		//use session
		s := sessionInfo{}
		s.At = time.Now()
		s.Id = oid
		return a.sess.Set(a.sessionName, s)
	}
	return nil
}
