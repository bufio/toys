// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package model provide some interfaces to support application development
on multiple database platforms.
*/
package model

import (
	"github.com/kidstuff/toys/util/errs"
	"sync"
)

var (
	drivers = make(map[string]Driver)
	mux     sync.Mutex
)

// Register makes a database driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, driver Driver) {
	mux.Lock()
	defer mux.Unlock()

	if driver == nil {
		panic("model: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("model: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

// Load return a Driver with name. An error return if the driver not exist
func Load(name string) (Driver, error) {
	mux.Lock()
	defer mux.Unlock()

	driver, ok := drivers[name]
	if !ok {
		return nil, errs.New("model: driver " + name + " not exist")
	}
	return driver, nil
}

// MusLoad return a Driver with name, it panic if the sriver not exist
func MustLoad(name string) Driver {
	driver, err := Load(name)
	if err != nil {
		panic("model: driver " + name + " not exist")
	}

	return driver
}

// Driver is the interface that ... do something...
type Driver interface {
	DecodeId(interface{}) (Identifier, error)
	ValidIdRep(interface{}) bool
	NewId() Identifier
}

type Identifier interface {
	Decode(interface{}) error
	Encode() string
	Valid() bool
}
