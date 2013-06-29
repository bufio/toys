// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package secure provide some convenient package to secure you web application.
*/
package secure

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

// RandomToken returns a random array of bytes
func RandomToken(l uint) []byte {
	t := make([]byte, l)
	if _, err := io.ReadFull(rand.Reader, t); err != nil {
		return nil
	}
	return t
}

func RandomString(l uint) string {
	// TODO: use a better solution rather than wasting some bytes
	l2 := l
	if d := l % 4; d != 0 {
		l2 = l2 + (4 - d)
	}
	return base64.StdEncoding.EncodeToString(RandomToken(l2 / 4 * 3))[:l]
}
