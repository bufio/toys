// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package secure

import (
	"crypto/rand"
	"io"
)

func RandomToken(l uint) []byte {
	t := make([]byte, l)
	if _, err := io.ReadFull(rand.Reader, t); err != nil {
		return nil
	}
	return t
}
