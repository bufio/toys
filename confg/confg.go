// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package confg

import (
	"github.com/kidstuff/toys/util/errs"
)

type Configurator interface {
	Load(path string) error
	Close() error
	Set(k string, v interface{})
	Get(k string) interface{}
	Del(k string)
}

var configurators = make(map[string]Configurator)

var (
	ErrConfiguratorNotFound = errs.New("confg: Configurator not found")
)

func Register(name string, configurator Configurator) {
	if configurator == nil {
		panic("confg: Register configurator is nil")
	}

	if _, dup := configurators[name]; dup {
		panic("confg: Register called twice for " + name)
	}

	configurators[name] = configurator
}

func Open(name, path string) (Configurator, error) {
	config, ok := configurators[name]
	if !ok {
		return nil, ErrConfiguratorNotFound
	}

	err := config.Load(path)
	if err != nil {
		return nil, errs.Err(err, "confg: cannot Open")
	}

	return config, nil
}
