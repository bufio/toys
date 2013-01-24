// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package membership

import (
	"log"
)

type Notificater interface {
	AccountAdded(email string, app bool) error
	PasswordChanged(email string) error
}

type SimpleNotificater struct{}

func NewSimpleNotificater() *SimpleNotificater {
	return &SimpleNotificater{}
}

func (n *SimpleNotificater) AccountAdded(email string, app bool) error {
	log.Printf("%s %t\n", email, app)
	return nil
}

func (n *SimpleNotificater) PasswordChanged(email string) error {
	log.Printf("%s\n", email)
	return nil
}
