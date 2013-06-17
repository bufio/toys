// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package model provide some interfaces to support application development
on multiple database platforms.
*/
package model

var drivers = make(map[string]Driver)

// Register makes a database driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, driver Driver) {
	if driver == nil {
		panic("model: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("model: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

type Driver interface {
	DecodeId(interface{}) (Identifier, error)
	ValidIdStr(interface{}) bool
	NewId() Identifier
}

type Identifier interface {
	Decode(interface{}) error
	Encode() string
	Valid() bool
}
