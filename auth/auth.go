// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"time"
)

type User interface {
	SetPassword(password string)
}

type Authenticater interface {
	SetOnlineTime(d time.Duration)
	GetOnlineTime() time.Duration
	CreateUser(email, password string) error
	DeleteUser(email string) error
	FindUser(id string) (User, error)
	FindUserByEmail(email string) (User, error)
	FindAllUser(offset string, limit int) ([]User, error)
	FindUserOnline(offset string, limit) ([]User, error)
	CountUserOnline() int
	ValidateUser(email, password string) bool
}
