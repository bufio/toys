// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package confg

import (
	"encoding/xml"
	"errors"
	"os"
)

type XmlSetting struct {
	Key   string
	Value interface{}
}

type XmlConfg struct {
	settings []XmlSetting
	file     *os.File
}

func NewXmlConfg(path string) (*XmlConfg, error) {
	f, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, errors.New("confg: cannot open congfig file")
	}

	var settings []XmlSetting
	dec := xml.NewDecoder(f)
	err = dec.Decode(settings)
	if err != nil {
		return nil, errors.New("confg: cannot decode stored data")
	}

	x := &XmlConfg{}
	x.file = f
	x.settings = settings
	return x, nil
}

func (x *XmlConfg) Close() {
	x.file.Close()
}

func (x *XmlConfg) Set(k string, v interface{}) {
	found := false
	for _, setting := range x.settings {
		if setting.Key == k {
			setting.Value = v
			found = true
		}
	}
	if !found {
		return
	}
}
