// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package confg

type Configurator interface {
	Set(k string, v interface{})
	Get(k string) interface{}
	Del(k string)
}
