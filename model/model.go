// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package model provide some interfaces to support application development
on multiple database platforms.
*/
package model

type Identifier interface {
	Decode(interface{}) error
	Encode() string
	Valid() bool
}
