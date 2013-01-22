// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package membership

type FormatChecker interface {
	PasswordValidate(string) bool
	EmailValidate(string) bool
}
